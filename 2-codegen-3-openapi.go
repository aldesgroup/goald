// ------------------------------------------------------------------------------------------------
// Here is the code used for generating the VMAP (value mapper) files
// ------------------------------------------------------------------------------------------------
package goald

import (
	"fmt"
	"log/slog"
	"net/http"
	"path"
	"strconv"

	"strings"
	"time"

	core "github.com/aldesgroup/corego"
	"github.com/aldesgroup/goald/features/utils"
	"github.com/getkin/kin-openapi/openapi3"
	"go.yaml.in/yaml/v3"
)

// ------------------------------------------------------------------------------------------------
// Main function with triggering logic
// ------------------------------------------------------------------------------------------------

const openAPIHTMLTemplate = `<!doctype html>
<html>
  <head>
    <title>API Reference</title>
    <meta charset="utf-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1" />
  </head>
  <body>
    <script id="api-reference" type="application/yaml">
%s
    </script>
    <script src="https://cdn.jsdelivr.net/npm/@scalar/api-reference"></script>
  </body>
</html>`

// Main function responsible for generating the API doc, eg data/api-doc.yaml (or custom path)
func (thisServer *server) generateOpenAPIDoc(srcdirs []string, docpath string, regen bool, servers string) {
	// checking the docpath
	core.PanicMsgIf(!strings.HasSuffix(docpath, ".yaml"), "The API doc path should end with .yaml")

	// first reason to generate: the regen is forced
	doRegen := regen

	// second reason: the API doc does not exist yet
	if !core.FileExists(docpath) {
		doRegen = true
	} else {
		// it exists, but is it recent enough? let's find out
		docModified := core.EnsureModTime(docpath)

		// let's check if a business object involved web exchanges has changed, or the endpoints code has changed
		doRegen = doRegen || isWebModelsChanged(docModified) || isWebCodeChanged(srcdirs, docModified)
	}

	// let's do it if we must
	if doRegen {

		// generating the doc content
		doc, errGen := generateOpenAPIYAML(servers)
		core.PanicMsgIfErr(errGen, "Error while generating the API doc in YAML")

		// writing it out
		core.WriteBytesToFile(docpath, doc)
		slog.Info("New version for: " + docpath)

		// also writing the HTML version
		core.WriteStringToFile(strings.Replace(docpath, ".yaml", ".html", 1), openAPIHTMLTemplate, string(doc))
	}
}

// Checking if at least 1 business object involved in an endpoint has changed
func isWebModelsChanged(docModified time.Time) bool {
	// going over all the endpoints
	for _, ep := range restRegistry.endpoints {
		if classRegistry.items[ep.getResourceClass()].getLastBOMod().After(docModified) {
			slog.Info(fmt.Sprintf("Output model '%s' for endpoint '%s %s' has changed!", ep.getResourceClass(), ep.getMethod(), ep.getLabel()))
			return true
		}
		if inputClass := ep.getInputOrParamsClass(); inputClass != "" && classRegistry.items[inputClass].getLastBOMod().After(docModified) {
			slog.Info(fmt.Sprintf("Input model '%s' for endpoint '%s %s' has changed!", ep.getResourceClass(), ep.getMethod(), ep.getLabel()))
			return true
		}
	}

	return false
}

// Checks in all the given source directories if some web code has changed and is more recent than the given date
func isWebCodeChanged(srcdirs []string, docModified time.Time) bool {
	for _, codedir := range srcdirs {
		if core.DirExists(codedir) {
			if checkWebCodeChanged(codedir, docModified) {
				return true
			}
		}
	}

	return false
}

// Checks in a the given source directory has some web code has changed more recent than the given date
func checkWebCodeChanged(codedir string, docModified time.Time) bool {
	for _, entry := range core.EnsureReadDir(codedir) {
		if entry.IsDir() {
			if checkWebCodeChanged(path.Join(codedir, entry.Name()), docModified) {
				return true
			}
		} else {
			if strings.HasSuffix(entry.Name(), "--web.go") {
				if filename := path.Join(codedir, entry.Name()); core.EnsureModTime(filename).After(docModified) {
					slog.Info("Endpoints might have changed with the latest change to this file: " + filename)
					return true
				}
			}
		}
	}

	return false
}

// ------------------------------------------------------------------------------------------------
// Utils - building an actual Open API doc from our declared endpoints
// ------------------------------------------------------------------------------------------------

// This builds the content of the doc, based on the config, and the coded endpoints
func generateOpenAPIYAML(servers string) ([]byte, error) {
	// the root of the doc
	doc := &openapi3.T{
		OpenAPI: "3.0.3",
		Info: &openapi3.Info{
			Title:       configObj.base().AppName,
			Version:     "DEVELOPMENT_VERSION", // this should always be replaced right before releasing a new version
			Description: configObj.base().AppDesc,
		},
		Paths: openapi3.NewPaths(),
		Components: &openapi3.Components{
			Schemas: openapi3.Schemas{},
		},
		Servers: openapi3.Servers{
			&openapi3.Server{
				URL:         "http://localhost:" + strconv.Itoa(configObj.base().Port),
				Description: "The local development environment",
			},
		},
	}

	// adding the remote servers
	for _, server := range strings.Split(servers, "|") {
		if parts := strings.SplitN(server, ":", 2); len(parts) == 2 {
			doc.Servers = append(doc.Servers, &openapi3.Server{
				URL:         parts[1],
				Description: fmt.Sprintf("The %s environment", parts[0]),
			})
		}
	}

	// adding each operations
	for _, ep := range getSortedEndpointList() {
		if err := addEndpointToDoc(doc, ep); err != nil {
			return nil, err
		}
	}

	// adding the tags, based on the groups of the endpoints
	for _, group := range core.GetSortedValues(allEndpointGroups) {
		doc.Tags = append(doc.Tags, &openapi3.Tag{
			Name:        group.Name,
			Description: group.Description,
		})
	}

	// marshaling it
	return yaml.Marshal(doc)
}

// ------------------------------------------------------------------------------------------------
// Utils - working at the operation level
// ------------------------------------------------------------------------------------------------

// This is used to add an operation to a doc
func addEndpointToDoc(doc *openapi3.T, ep iEndpoint) error {
	// the complete operation's path, using OpenAPI {param} syntax instead of httprouter :param
	basePath := ep.getOperationPath(false)
	openapiPath := basePath
	operationID := strings.ToLower(ep.getMethod()) + basePath
	if idProp := ep.getIDProp(); idProp != nil {
		openapiPath += "/{" + idProp.getName() + "}"
		operationID += "/" + idProp.getName()
	}

	pathItem := doc.Paths.Find(openapiPath)
	if pathItem == nil {
		pathItem = &openapi3.PathItem{}
		doc.Paths.Set(openapiPath, pathItem)
	}

	if ep.getGroup() == nil {
		return fmt.Errorf("Endpoint '%s' does not have a group", ep.getPathAsString())
	}

	op := &openapi3.Operation{
		Summary:     ep.getLabel(),
		Description: ep.getDescription(),
		Responses:   &openapi3.Responses{},
		OperationID: operationID,
		Tags:        []string{ep.getGroup().Name},
	}

	// ---------- PATH PARAM ----------
	if idProp := ep.getIDProp(); idProp != nil {
		op.Parameters = append(op.Parameters, &openapi3.ParameterRef{
			Value: &openapi3.Parameter{
				Name:        idProp.getName(),
				In:          "path",
				Required:    true,
				Schema:      schemaFromPrimitiveType(idProp, false),
				Description: "URL path parameter: " + idProp.getTag("desc"),
			},
		})
	}

	// ---------- INPUT ----------
	inClass := ep.getInputOrParamsClass()
	if inClass != "" {

		if ep.isBodyInputRequired() {
			// === BODY INPUT ===
			ref, err := getSchemaRef(doc, inClass)
			if err != nil {
				return err
			}

			var description string
			if ep.isMultipleInput() {
				description = fmt.Sprintf("An array of %s objects", inClass)
			} else {
				description = fmt.Sprintf("A %s object", inClass)
			}

			op.RequestBody = &openapi3.RequestBodyRef{
				Value: &openapi3.RequestBody{
					Required: true,
					Content: openapi3.Content{
						"application/json": &openapi3.MediaType{
							Schema: ref,
						},
					},
					Description: description,
				},
			}

		} else {
			// === URL PARAMS INPUT ===
			params, err := paramsFromClass(inClass, openapiPath)
			if err != nil {
				return err
			}
			op.Parameters = append(op.Parameters, params...)
		}
	}

	// ---------- OUTPUT ----------
	outClass := ep.getResourceClass()
	if outClass != "" {
		ref, err := getSchemaRef(doc, outClass)
		if err != nil {
			return err
		}

		op.Responses.Set("200", &openapi3.ResponseRef{
			Value: &openapi3.Response{
				Description: strPtr("OK"),
				Content: openapi3.Content{
					"application/json": &openapi3.MediaType{
						Schema: ref,
					},
				},
			},
		})
	}

	// ---------- METHOD BINDING ----------
	switch ep.getMethod() {
	case http.MethodGet:
		pathItem.Get = op
	case http.MethodPost:
		pathItem.Post = op
	case http.MethodPut:
		pathItem.Put = op
	case http.MethodPatch:
		pathItem.Patch = op
	case http.MethodDelete:
		pathItem.Delete = op
	default:
		return fmt.Errorf("unsupported HTTP method %q", ep.getMethod())
	}

	return nil
}

// ------------------------------------------------------------------------------------------------
// Utils - schemas
// ------------------------------------------------------------------------------------------------

// we'll put all the schemas here, and avoid repeating ourselves
var schemaCache = map[className]*openapi3.SchemaRef{}

// getting a schema REF for a given class, initializing it if needed
func getSchemaRef(doc *openapi3.T, clsName className) (*openapi3.SchemaRef, error) {
	// fast returning if possible
	if clsName == "" {
		return nil, nil
	}
	if ref, ok := schemaCache[clsName]; ok {
		return ref, nil
	}

	// early caching of a new schema REF for the given class, to avoid cycles
	ref := &openapi3.SchemaRef{
		Ref: "#/components/schemas/" + string(clsName),
	}
	schemaCache[clsName] = ref

	// getting the associated model
	model := modelForName(clsName)
	if model == nil {
		return nil, fmt.Errorf("No model associated with BO class '%s'", clsName)
	}

	// adding a schema (not juste a REF) for the given class to the doc's components
	doc.Components.Schemas[string(clsName)] = &openapi3.SchemaRef{Value: schemaFromModel(doc, model)}

	return ref, nil
}

// building a schema for the given business object model
func schemaFromModel(doc *openapi3.T, model IBusinessObjectModel) *openapi3.Schema {

	// new schema
	schema := &openapi3.Schema{
		Type:        &openapi3.Types{"object"},
		Properties:  openapi3.Schemas{},
		Description: model.base().description,
	}

	// going over the basic properties, i.e. the fields
	for _, field := range core.GetSortedValues(model.base().fields) {
		fieldJSONName := core.PascalToCamel(field.getName())

		var prop *openapi3.SchemaRef = schemaFromPrimitiveType(field, true)

		// linking this field to the schema
		schema.Properties[fieldJSONName] = prop

		// gathering the required
		if field.isMandatoryInput() {
			schema.Required = append(schema.Required, fieldJSONName)
		} else if field.isPureOutput() {
			prop.Value.ReadOnly = true
		}
	}

	// going over the relationships, i.e. the object-type properties
	for _, relationship := range core.GetSortedValues(model.base().relationships) {
		relationshipJSONName := core.PascalToCamel(relationship.getName())

		if !relationship.polymorphic || len(relationship.targets) == 1 {

			// getting the schema REF for the relationship target
			prop, errRef := getSchemaRef(doc, relationship.targets[0].base().name)
			core.PanicMsgIfErr(errRef, "Error while getting schema ref for relationship '%s#%s'",
				model.base().name, relationship.getName())

			// linking this relationship to the schema
			if relationship.isPureOutput() {
				schema.Properties[relationshipJSONName] = &openapi3.SchemaRef{
					Value: &openapi3.Schema{
						ReadOnly: true,
						AllOf:    openapi3.SchemaRefs{prop},
					},
				}
			} else {
				schema.Properties[relationshipJSONName] = prop
			}

			// gathering the required
			if relationship.isMandatoryInput() {
				schema.Required = append(schema.Required, relationshipJSONName)
			}
		} else {
			schema.Properties[relationshipJSONName] = &openapi3.SchemaRef{
				Value: &openapi3.Schema{
					Type:  &openapi3.Types{"object"},
					OneOf: []*openapi3.SchemaRef{},
				},
			}

			if len(relationship.targets) == 0 {
				core.PanicMsg("Polymorphic relationship '%s#%s' does not have any target model", model.base().name, relationship.getName())
			}

			for _, target := range relationship.targets {
				prop, errRef := getSchemaRef(doc, target.base().name)
				core.PanicMsgIfErr(errRef, "Error while getting schema ref for relationship '%s#%s'",
					model.base().name, relationship.getName())
				schema.Properties[relationshipJSONName].Value.OneOf = append(schema.Properties[relationshipJSONName].Value.OneOf, prop)
			}

			if relationship.isMandatoryInput() {
				schema.Required = append(schema.Required, relationshipJSONName)
			} else if relationship.isPureOutput() {
				schema.Properties[relationshipJSONName].Value.ReadOnly = true
			}
		}
	}

	return schema
}

// building URL parameters from the given business object model that's associated with a URLQueryParams-derived BO
func paramsFromClass(clsName className, path string) (openapi3.Parameters, error) {
	model := modelForName(clsName)
	if model == nil {
		return nil, fmt.Errorf("No model associated with URL Query Params class '%s'", clsName)
	}

	// pathVars := extractPathVars(path)
	var out openapi3.Parameters

	for _, field := range core.GetSortedValues(model.base().fields) {
		schema := schemaFromPrimitiveType(field, true)
		parameter := &openapi3.Parameter{
			Name:        field.getName(),
			In:          "query",
			Required:    field.isMandatoryInput(),
			Schema:      schema,
			Description: "URL query Parameter: " + schema.Value.Description,
		}

		out = append(out, &openapi3.ParameterRef{Value: parameter})
	}

	return out, nil
}

func withDescription(schema *openapi3.Schema, description string) *openapi3.Schema {
	if schema != nil && description != "" {
		schema.Description = description
	}
	return schema
}

func schemaFromPrimitiveType(field IField, addDesc bool) *openapi3.SchemaRef {

	var description string
	if addDesc {
		if field.getName() == "ID" {
			description = "the unique identifier of the " + string(field.ownerModel().base().name)
		} else {
			description = field.getTag("desc")
		}
	}

	switch field.getTypeFamily() {

	case utils.TypeFamilyBOOL:
		return &openapi3.SchemaRef{Value: withDescription(openapi3.NewBoolSchema(), description)}

	case utils.TypeFamilySTRING:
		return &openapi3.SchemaRef{Value: withDescription(openapi3.NewStringSchema(), description)}

	case utils.TypeFamilyINT:
		return &openapi3.SchemaRef{Value: withDescription(openapi3.NewInt32Schema(), description)}

	case utils.TypeFamilyBIGINT:
		return &openapi3.SchemaRef{Value: withDescription(openapi3.NewInt64Schema(), description)}

	case utils.TypeFamilyREAL, utils.TypeFamilyDOUBLE:
		return &openapi3.SchemaRef{Value: withDescription(openapi3.NewFloat64Schema(), description)}

	// case utils.TypeFamilyDATE:               "date", // TODO

	case utils.TypeFamilyENUM:
		// instantiating an instance of the owner of this field
		enumOwner := getClass(field.ownerModel()).NewObject()

		// this owner has a zero-value for this field, which is enough for us to do the rest
		enumVal := utils.ValueOf(enumOwner).GetFieldValue(field.getName())

		// controlling we do have an enum
		if enum, ok := enumVal.(IEnum); ok {
			// getting all this enum's values + corresponding labels
			values := enum.Values()

			// getting only the integer values, sorted to allow diffing
			enumVals := core.GetSortedKeys(values)

			// we'll display the value-label pairs as a markdown list
			var lines []string
			for _, val := range enumVals {
				lines = append(lines, fmt.Sprintf("- `%d` = %s", val, values[val]))
			}

			return &openapi3.SchemaRef{
				Value: &openapi3.Schema{
					Type:        &openapi3.Types{"integer"},
					Format:      "int32",
					Enum:        core.ToAnySlice(enumVals),
					Description: field.getTag("desc") + "; labels associated with the possible values: \n" + strings.Join(lines, "\n"),
				}}

		} else {
			// should never happen
			core.PanicMsg("It seems field '%s' is not a proper enum (does not implement goald.IEnum)", field.getName())
			return nil
		}

	// case "array":
	// 	return &openapi3.SchemaRef{Value: &openapi3.Schema{
	// 		Type:  &openapi3.Types{"array"},
	// 		Items: schemaFromPrimitiveType(*t.ArrayItem),
	// 	}}

	default:
		panic("Unhandled type in Open API doc generation: " + field.getTypeFamily().String())
	}
}

func strPtr(s string) *string { return &s }

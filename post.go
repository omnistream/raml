package raml

import (
	//"log"
	"strings"
)

func (d *APIDefinition) PostProcessing() {
	//postProcess resource
	for k, resource := range d.Resources {
		resource.URI = k
		d.PostResource(resource)
	}
}

func (d *APIDefinition) PostResource(resource *Resource) *Resource {
	//log.Println("PostResource ", resource.URI)
	if resource.Get != nil {
		resource.Get = d.PostResourceGet(resource.Get)
	}
	if resource.Post != nil {
		resource.Post = d.PostResourcePost(resource.Post)
	}
	/*
		if resource.Put != nil {
			fmt.Println("Is Put")
			method = "PUT"
		}
		if resource.Delete != nil {
			fmt.Println("Is Delete")
			method = "DELETE"
		}
		if resource.Head != nil {
			fmt.Println("Is Head")
			method = "HEAD"
		}
		if resource.Patch != nil {
			fmt.Println("Is Patch")
			method = "PATCH"
		}
	*/

	resource.SecuredBy = d.PostSecuredBy(resource.TempSecuredBy)
	if resource.Nested != nil {
		for k, resourceNested := range resource.Nested {
			//log.Printf("nesting for %s \n", k)
			resourceNested.URI = resource.URI + k
			d.PostResource(resourceNested)
			//			fmt.Printf("calling %s\n", resource.URI)
		}
	}

	return resource
}

func (d *APIDefinition) PostResourceGet(method *Method) *Method {
	//log.Println("PostResourceGet")
	method.Bodies = d.PostResourceBodies(method.TempBodies)
	d.PostResourceResponses(method.Responses)
	method.SecuredBy = d.PostSecuredBy(method.TempSecuredBy)
	return method
}

func (d *APIDefinition) PostResourcePost(method *Method) *Method {
	//log.Println("PostResourcePost")
	method.Bodies = d.PostResourceBodies(method.TempBodies)
	d.PostResourceResponses(method.Responses)
	method.SecuredBy = d.PostSecuredBy(method.TempSecuredBy)
	return method
}

func (d *APIDefinition) PostResourceBodies(tempBodies *TempBodies) *Bodies {
	//log.Println("PostResourceBodies")
	bodies := &Bodies{}
	if tempBodies != nil {
		if tempBodies.ForMIMEType != nil {
			for contentType, body := range tempBodies.ForMIMEType {
				d.PostResourceBody(contentType, body, bodies)
				//log.Println("bodies2:", bodies)
			}
		} else {
			//				mediaType: contentType,
			/*
				bodies.Add(Body{
					Schema:  tempBodies.DefaultSchema,
					Example: tempBodies.DefaultExample,
				})
			*/

			/*
				PostResourceBody(contentType, tempBodies, bodies)
				log.Println("bodies3:", bodies)
			*/
		}
	}
	return bodies
}

func (d *APIDefinition) PostResourceBody(contentType string, body Body, bodies *Bodies) {
	//log.Println("PostResourceBody")
	var (
		schema  string
		example string
	)

	if strings.Index(body.Schema, "{") == -1 {
		tmp := findInSlice(d.Schemas, body.Schema)
		schema = strings.Replace(tmp, "\n", " ", -1)
	} else {
		schema = strings.Replace(body.Schema, "\n", " ", -1)
	}
	//log.Println("Index", body.Schema, " ", strings.Index(body.Schema, "{"))

	if strings.Index(body.Example, "{") == -1 {
		tmp := findInSlice(d.Schemas, body.Example)
		example = strings.Replace(tmp, "\n", " ", -1)
	} else {
		//		strings.Replace(s,"\n","<br>",-1)
		example = strings.Replace(body.Example, "\n", " ", -1)
	}
	//log.Println("Index", body.Example, " ", strings.Index(body.Example, "{"))

	bodies.Add(Body{
		MediaType: contentType,
		Schema:    schema,
		Example:   example,
	})
	//log.Println("bodies:", bodies)
	/*
		bodies = append(bodies, Body{
			mediaType: contentType,
			Schema:    body.Schema,
			Example:   body.Example,
		})
	*/
}

func (d *APIDefinition) PostResourceResponses(responses map[HTTPCode]*Response) {
	//log.Println("PostResourceResponses")
	for _, response := range responses {
		response.Bodies = d.PostResourceBodies(response.TempBodies)
	}
}

func (d *APIDefinition) PostResourceResponse(code HTTPCode, response Response) {
	//log.Println("PostResourceResponse")
}

func (d *APIDefinition) PostSecuredBy(temp []interface{}) []*DefinitionChoice {
	var result []*DefinitionChoice
	//log.Printf("PostSecuredBy %#v \n", temp)
	for _, v := range temp {
		//log.Printf("PostSecuredBy2 %#v \n", v)
		if t, ok := v.(string); ok {
			res := &DefinitionChoice{
				Name: t,
			}
			result = append(result, res)
		}
	}
	return result
}
func findInSlice(slice []map[string]string, name string) string {
	for _, v := range slice {
		if v[name] != "" {
			return v[name]
		}
	}
	return ""
}

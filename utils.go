package raml

import (
//	"log"
)

func (d *APIDefinition) GetSecurityHeaders(uri, method string) *map[string]string {
	var (
		res  *Resource
		meth *Method
	)
	headers := &map[string]string{}
	for k, r := range d.Resources {
		r.URI = k
		res, meth = parseResource(uri, method, r)
		if res != nil {
			break
		}
	}
	if res != nil {
		//		log.Println("we found ", res.URI)
		if meth.SecuredBy != nil {
			for _, sec := range meth.SecuredBy {
				//				log.Printf("secured by %#v \n", sec)
				security := findSecurity(d, sec.Name)
				if security != nil {
					if security.DescribedBy != nil {
						if security.DescribedBy.Headers != nil {
							for k, v := range security.DescribedBy.Headers {
								//								log.Printf("header %#v %#v\n", k, v)
								(*headers)[string(k)] = v.Example
							}
						}
					}
				}
			}
		}
	}
	return headers
}

func parseResource(uri, method string, resource *Resource) (*Resource, *Method) {
	if uri == resource.URI {
		if method == "GET" && resource.Get != nil {
			return resource, resource.Get
		}
		if method == "POST" && resource.Post != nil {
			return resource, resource.Post
		}
		if method == "DELETE" && resource.Delete != nil {
			return resource, resource.Delete
		}
		if method == "HEAD" && resource.Head != nil {
			return resource, resource.Head
		}
		if method == "PUT" && resource.Put != nil {
			return resource, resource.Put
		}
		if method == "PATCH" && resource.Patch != nil {
			return resource, resource.Patch
		}
	}
	if resource.Nested != nil {
		for k, resourceNested := range resource.Nested {
			resourceNested.URI = resource.URI + k
			return parseResource(uri, method, resourceNested)
		}
	}
	return nil, nil
}

func findSecurity(api *APIDefinition, name string) *SecurityScheme {
	for _, sec := range api.SecuritySchemes {
		if sec[name] != nil {
			return sec[name]
		}
	}
	return nil
}

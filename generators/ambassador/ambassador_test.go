package ambassador

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/kubeshop/kusk/options"
	"github.com/kubeshop/kusk/spec"
)

type testCase struct {
	name    string
	options options.Options
	spec    string
	res     string
}

func TestAmbassador(t *testing.T) {
	var gen Generator

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			r := require.New(t)

			spec, err := spec.Parse([]byte(testCase.spec))
			r.NoError(err, "failed to parse spec")

			mappings, err := gen.Generate(&testCase.options, spec)
			r.NoError(err)
			r.Equal(testCase.res, mappings)
		})
	}
}

var testCases = []testCase{
	{
		name: "basic",
		options: options.Options{
			Namespace: "default",
			Service: options.ServiceOptions{
				Namespace: "default",
				Name:      "petstore",
			},
			Path: options.PathOptions{
				Base:       "",
				TrimPrefix: "",
				Split:      true,
			},
		},
		spec: `
openapi: 3.0.2
info:
  title: Swagger Petstore - OpenAPI 3.0
  version: 1.0.5
x-kusk:
  namespace: notdefault
  service:
    name: petstore
paths:
  "/pet":
    put:
      operationId: updatePet
      responses:
        '200':
          description: Successful operation
`,
		res: `
---
apiVersion: getambassador.io/v2
kind: Mapping
metadata:
  name: petstore-updatepet
  namespace: default
spec:
  prefix: "/pet"
  method: PUT
  service: petstore.default:80
  rewrite: ""
`,
	},
	{
		name: "basic-json",
		options: options.Options{
			Namespace: "default",
			Service: options.ServiceOptions{
				Namespace: "default",
				Name:      "petstore",
			},
			Path: options.PathOptions{
				Base:       "",
				TrimPrefix: "",
				Split:      true,
			},
		},
		spec: `
{
  "openapi": "3.0.2",
  "info": {
    "title": "Swagger Petstore - OpenAPI 3.0",
    "version": "1.0.5"
  },
  "paths": {
    "/pet": {
      "put": {
        "operationId": "updatePet",
        "responses": {
          "200": {
            "description": "Successful operation"
          }
        }
      }
    }
  }
}
`,
		res: `
---
apiVersion: getambassador.io/v2
kind: Mapping
metadata:
  name: petstore-updatepet
  namespace: default
spec:
  prefix: "/pet"
  method: PUT
  service: petstore.default:80
  rewrite: ""
`,
	},
	{
		name: "basic-namespace",
		options: options.Options{
			Namespace: "amb",
			Service: options.ServiceOptions{
				Namespace: "default",
				Name:      "petstore",
			},
			Path: options.PathOptions{
				Base:       "",
				TrimPrefix: "",
				Split:      true,
			},
		},
		spec: `
openapi: 3.0.2
info:
  title: Swagger Petstore - OpenAPI 3.0
  version: 1.0.5
paths:
  "/pet":
    put:
      operationId: updatePet
      responses:
        '200':
          description: Successful operation
`,
		res: `
---
apiVersion: getambassador.io/v2
kind: Mapping
metadata:
  name: petstore-updatepet
  namespace: amb
spec:
  prefix: "/pet"
  method: PUT
  service: petstore.default:80
  rewrite: ""
`,
	},
	{
		name: "parameter",
		options: options.Options{
			Namespace: "default",
			Service: options.ServiceOptions{
				Namespace: "default",
				Name:      "petstore",
			},
			Path: options.PathOptions{
				Base:       "",
				TrimPrefix: "",
				Split:      true,
			},
		},
		spec: `
openapi: 3.0.2
info:
  title: Swagger Petstore - OpenAPI 3.0
  version: 1.0.5
paths:
  "/pet/{petId}/uploadImage":
    post:
      operationId: uploadFile
      parameters:
        - name: petId
          in: path
          description: ID of pet to update
          required: true
          schema:
            type: integer
            format: int64
      responses:
        '200':
          description: Successful operation
`,
		res: `
---
apiVersion: getambassador.io/v2
kind: Mapping
metadata:
  name: petstore-uploadfile
  namespace: default
spec:
  prefix: "/pet/([a-zA-Z0-9]*)/uploadImage"
  prefix_regex: true
  method: POST
  service: petstore.default:80
  rewrite: ""
`,
	},
	{
		name: "empty-operationId",
		options: options.Options{
			Namespace: "default",
			Service: options.ServiceOptions{
				Namespace: "default",
				Name:      "petstore",
			},
			Path: options.PathOptions{
				Base:       "",
				TrimPrefix: "",
				Split:      true,
			},
		},
		spec: `
openapi: 3.0.2
info:
  title: Swagger Petstore - OpenAPI 3.0
  version: 1.0.5
paths:
  "/pet/{petId}/uploadImage":
    post:
      parameters:
        - name: petId
          in: path
          description: ID of pet to update
          required: true
          schema:
            type: integer
            format: int64
      responses:
        '200':
          description: Successful operation
`,
		res: `
---
apiVersion: getambassador.io/v2
kind: Mapping
metadata:
  name: petstore-postpetpetiduploadimage
  namespace: default
spec:
  prefix: "/pet/([a-zA-Z0-9]*)/uploadImage"
  prefix_regex: true
  method: POST
  service: petstore.default:80
  rewrite: ""
`,
	},
	{
		name: "basepath",
		options: options.Options{
			Namespace: "default",
			Service: options.ServiceOptions{
				Namespace: "default",
				Name:      "petstore",
			},
			Path: options.PathOptions{
				Base:       "/api/v3",
				TrimPrefix: "",
				Split:      true,
			},
		},
		spec: `
openapi: 3.0.2
info:
  title: Swagger Petstore - OpenAPI 3.0
  version: 1.0.5
paths:
  "/pet/{petId}/uploadImage":
    post:
      parameters:
        - name: petId
          in: path
          description: ID of pet to update
          required: true
          schema:
            type: integer
            format: int64
      responses:
        '200':
          description: Successful operation
`,
		res: `
---
apiVersion: getambassador.io/v2
kind: Mapping
metadata:
  name: petstore-postpetpetiduploadimage
  namespace: default
spec:
  prefix: "/api/v3/pet/([a-zA-Z0-9]*)/uploadImage"
  prefix_regex: true
  method: POST
  service: petstore.default:80
  rewrite: ""
`,
	},
	{
		name: "basepath-rootonly",
		options: options.Options{
			Namespace: "default",
			Service: options.ServiceOptions{
				Namespace: "default",
				Name:      "petstore",
			},
			Path: options.PathOptions{
				Base:       "/api/v3",
				TrimPrefix: "",
				Split:      false,
			},
		},
		spec: `
openapi: 3.0.2
info:
  title: Swagger Petstore - OpenAPI 3.0
  version: 1.0.5
paths:
  "/pet":
    put:
      operationId: updatePet
      responses:
        '200':
          description: Successful operation
  "/pet/{petId}/uploadImage":
    post:
      operationId: uploadFile
      parameters:
        - name: petId
          in: path
          description: ID of pet to update
          required: true
          schema:
            type: integer
            format: int64
      responses:
        '200':
          description: Successful operation`,
		res: `
---
apiVersion: getambassador.io/v2
kind: Mapping
metadata:
  name: petstore
  namespace: default
spec:
  prefix: "/api/v3"
  service: petstore.default:80
  rewrite: ""
`,
	},
	{
		name: "basepath-trimprefix",
		options: options.Options{
			Namespace: "default",
			Service: options.ServiceOptions{
				Namespace: "default",
				Name:      "petstore",
			},
			Path: options.PathOptions{
				Base:       "/petstore/api/v3",
				TrimPrefix: "/petstore",
				Split:      true,
			},
		},
		spec: `
openapi: 3.0.2
info:
  title: Swagger Petstore - OpenAPI 3.0
  version: 1.0.5
paths:
  "/pet/{petId}/uploadImage":
    post:
      parameters:
        - name: petId
          in: path
          description: ID of pet to update
          required: true
          schema:
            type: integer
            format: int64
      responses:
        '200':
          description: Successful operation
`,
		res: `
---
apiVersion: getambassador.io/v2
kind: Mapping
metadata:
  name: petstore-postpetpetiduploadimage
  namespace: default
spec:
  prefix: "/petstore/api/v3/pet/([a-zA-Z0-9]*)/uploadImage"
  prefix_regex: true
  method: POST
  service: petstore.default:80
  regex_rewrite:
    pattern: '/petstore(.*)'
    substitution: '\1'
`,
	},
	{
		name: "swagger-yaml",
		options: options.Options{
			Namespace: "default",
			Service: options.ServiceOptions{
				Namespace: "default",
				Name:      "petstore",
			},
			Path: options.PathOptions{
				Base:       "",
				TrimPrefix: "",
				Split:      true,
			},
		},
		spec: `
swagger: "2.0"
info:
  version: 1.0.0
  title: Swagger Petstore
basePath: /v1
paths:
  /pets:
    get:
      summary: List all pets
      operationId: listPets
      parameters:
        - name: limit
          in: query
          required: false
          type: integer
          format: int32
      responses:
        "200":
          description: A paged array of pets
          schema:
            $ref: '#/definitions/Pets'
        default:
          description: unexpected error
          schema:
            $ref: '#/definitions/Error'
    post:
      summary: Create a pet
      operationId: createPets
      responses:
        "201":
          description: Null response
        default:
          description: unexpected error
          schema:
            $ref: '#/definitions/Error'
  /pets/{petId}:
    get:
      operationId: showPetById
      parameters:
        - name: petId
          in: path
          required: true
          type: string
      responses:
        "200":
          description: Expected response to a valid request
          schema:
            $ref: '#/definitions/Pets'
        default:
          description: unexpected error
          schema:
            $ref: '#/definitions/Error'
definitions:
  Pet:
    type: "object"
    required:
      - id
      - name
    properties:
      id:
        type: integer
        format: int64
      name:
        type: string
      tag:
        type: string
  Pets:
    type: array
    items:
      $ref: '#/definitions/Pet'
  Error:
    type: "object"
    required:
      - code
      - message
    properties:
      code:
        type: integer
        format: int32
      message:
        type: string
`,
		res: `
---
apiVersion: getambassador.io/v2
kind: Mapping
metadata:
  name: petstore-createpets
  namespace: default
spec:
  prefix: "/pets"
  method: POST
  service: petstore.default:80
  rewrite: ""
---
apiVersion: getambassador.io/v2
kind: Mapping
metadata:
  name: petstore-listpets
  namespace: default
spec:
  prefix: "/pets"
  method: GET
  service: petstore.default:80
  rewrite: ""
---
apiVersion: getambassador.io/v2
kind: Mapping
metadata:
  name: petstore-showpetbyid
  namespace: default
spec:
  prefix: "/pets/([a-zA-Z0-9]*)"
  prefix_regex: true
  method: GET
  service: petstore.default:80
  rewrite: ""
`,
	},
	{
		name: "swagger-json",
		options: options.Options{
			Namespace: "default",
			Service: options.ServiceOptions{
				Namespace: "default",
				Name:      "petstore",
			},
			Path: options.PathOptions{
				Base:       "",
				TrimPrefix: "",
				Split:      true,
			},
		},
		spec: `
{
  "swagger": "2.0",
  "info": {
    "version": "1.0.0",
    "title": "Swagger Petstore"
  },
  "basePath": "/v1",
  "paths": {
    "/pets": {
      "get": {
        "summary": "List all pets",
        "operationId": "listPets",
        "parameters": [
          {
            "name": "limit",
            "in": "query",
            "required": false,
            "type": "integer",
            "format": "int32"
          }
        ],
        "responses": {
          "200": {
            "description": "A paged array of pets",
            "schema": {
              "$ref": "#/definitions/Pets"
            }
          },
          "default": {
            "description": "unexpected error",
            "schema": {
              "$ref": "#/definitions/Error"
            }
          }
        }
      },
      "post": {
        "summary": "Create a pet",
        "operationId": "createPets",
        "responses": {
          "201": {
            "description": "Null response"
          },
          "default": {
            "description": "unexpected error",
            "schema": {
              "$ref": "#/definitions/Error"
            }
          }
        }
      }
    },
    "/pets/{petId}": {
      "get": {
        "operationId": "showPetById",
        "parameters": [
          {
            "name": "petId",
            "in": "path",
            "required": true,
            "type": "string"
          }
        ],
        "responses": {
          "200": {
            "description": "Expected response to a valid request",
            "schema": {
              "$ref": "#/definitions/Pets"
            }
          },
          "default": {
            "description": "unexpected error",
            "schema": {
              "$ref": "#/definitions/Error"
            }
          }
        }
      }
    }
  },
  "definitions": {
    "Pet": {
      "type": "object",
      "required": [
        "id",
        "name"
      ],
      "properties": {
        "id": {
          "type": "integer",
          "format": "int64"
        },
        "name": {
          "type": "string"
        },
        "tag": {
          "type": "string"
        }
      }
    },
    "Pets": {
      "type": "array",
      "items": {
        "$ref": "#/definitions/Pet"
      }
    },
    "Error": {
      "type": "object",
      "required": [
        "code",
        "message"
      ],
      "properties": {
        "code": {
          "type": "integer",
          "format": "int32"
        },
        "message": {
          "type": "string"
        }
      }
    }
  }
}
`,
		res: `
---
apiVersion: getambassador.io/v2
kind: Mapping
metadata:
  name: petstore-createpets
  namespace: default
spec:
  prefix: "/pets"
  method: POST
  service: petstore.default:80
  rewrite: ""
---
apiVersion: getambassador.io/v2
kind: Mapping
metadata:
  name: petstore-listpets
  namespace: default
spec:
  prefix: "/pets"
  method: GET
  service: petstore.default:80
  rewrite: ""
---
apiVersion: getambassador.io/v2
kind: Mapping
metadata:
  name: petstore-showpetbyid
  namespace: default
spec:
  prefix: "/pets/([a-zA-Z0-9]*)"
  prefix_regex: true
  method: GET
  service: petstore.default:80
  rewrite: ""
`,
	},
	{
		name: "port specified",
		options: options.Options{
			Namespace: "default",
			Service: options.ServiceOptions{
				Namespace: "default",
				Name:      "petstore",
				Port:      443,
			},
			Path: options.PathOptions{
				Base:       "",
				TrimPrefix: "",
				Split:      true,
			},
		},
		spec: `
openapi: 3.0.2
info:
  title: Swagger Petstore - OpenAPI 3.0
  version: 1.0.5
paths:
  "/pet":
    put:
      operationId: updatePet
      responses:
        '200':
          description: Successful operation
`,
		res: `
---
apiVersion: getambassador.io/v2
kind: Mapping
metadata:
  name: petstore-updatepet
  namespace: default
spec:
  prefix: "/pet"
  method: PUT
  service: petstore.default:443
  rewrite: ""
`,
	},
	{
		name: "port 0 specified",
		options: options.Options{
			Namespace: "default",
			Service: options.ServiceOptions{
				Namespace: "default",
				Name:      "petstore",
				Port:      0,
			},
			Path: options.PathOptions{
				Base:       "",
				TrimPrefix: "",
				Split:      true,
			},
		},
		spec: `
openapi: 3.0.2
info:
  title: Swagger Petstore - OpenAPI 3.0
  version: 1.0.5
paths:
  "/pet":
    put:
      operationId: updatePet
      responses:
        '200':
          description: Successful operation
`,
		res: `
---
apiVersion: getambassador.io/v2
kind: Mapping
metadata:
  name: petstore-updatepet
  namespace: default
spec:
  prefix: "/pet"
  method: PUT
  service: petstore.default:80
  rewrite: ""
`,
	},
	{
		name: "path-disabled",
		options: options.Options{
			Namespace: "default",
			Service: options.ServiceOptions{
				Namespace: "default",
				Name:      "petstore",
			},
			PathSubOptions: map[string]options.SubOptions{
				"/pet": {
					Disabled: true,
				},
			},
		},
		spec: `
openapi: 3.0.2
info:
  title: Swagger Petstore - OpenAPI 3.0
  version: 1.0.5
paths:
  "/pet":
    x-kusk:
      disabled: true
    put:
      operationId: updatePet
      responses:
        '200':
          description: Successful operation
  "/pet/{petId}/uploadImage":
    post:
      operationId: uploadFile
      parameters:
        - name: petId
          in: path
          description: ID of pet to update
          required: true
          schema:
            type: integer
            format: int64
      responses:
        '200':
          description: Successful operation`,
		res: `
---
apiVersion: getambassador.io/v2
kind: Mapping
metadata:
  name: petstore-uploadfile
  namespace: default
spec:
  prefix: "/pet/([a-zA-Z0-9]*)/uploadImage"
  prefix_regex: true
  method: POST
  service: petstore.default:80
  rewrite: ""
`,
	},
	{
		name: "operation-disabled",
		options: options.Options{
			Namespace: "default",
			Service: options.ServiceOptions{
				Namespace: "default",
				Name:      "petstore",
			},
			OperationSubOptions: map[string]options.SubOptions{
				"PUT/pet": {
					Disabled: true,
				},
			},
		},
		spec: `
openapi: 3.0.2
info:
  title: Swagger Petstore - OpenAPI 3.0
  version: 1.0.5
paths:
  "/pet":
    put:
      x-kusk:
        disabled: true
      operationId: updatePet
      responses:
        '200':
          description: Successful operation
  "/pet/{petId}/uploadImage":
    post:
      operationId: uploadFile
      parameters:
        - name: petId
          in: path
          description: ID of pet to update
          required: true
          schema:
            type: integer
            format: int64
      responses:
        '200':
          description: Successful operation`,
		res: `
---
apiVersion: getambassador.io/v2
kind: Mapping
metadata:
  name: petstore-uploadfile
  namespace: default
spec:
  prefix: "/pet/([a-zA-Z0-9]*)/uploadImage"
  prefix_regex: true
  method: POST
  service: petstore.default:80
  rewrite: ""
`,
	},
	{
		name: "cors-global",
		options: options.Options{
			Namespace: "default",
			Service: options.ServiceOptions{
				Namespace: "default",
				Name:      "petstore",
			},
			Ingress: options.IngressOptions{
				CORS: options.CORSOptions{
					Origins:       []string{"http://foo.example", "http://bar.example"},
					Methods:       []string{"POST", "GET", "OPTIONS"},
					Headers:       []string{"Content-Type"},
					ExposeHeaders: []string{"X-Custom-Header", "X-Other-Custom-Header"},
					Credentials:   nil,
					MaxAge:        120,
				},
			},
		},
		spec: `
openapi: 3.0.2
info:
  title: Swagger Petstore - OpenAPI 3.0
  version: 1.0.5
paths:
  "/pet":
    put:
      operationId: updatePet
      responses:
        '200':
          description: Successful operation
  "/pet/{petId}/uploadImage":
    post:
      operationId: uploadFile
      parameters:
        - name: petId
          in: path
          description: ID of pet to update
          required: true
          schema:
            type: integer
            format: int64
      responses:
        '200':
          description: Successful operation`,
		res: `
---
apiVersion: getambassador.io/v2
kind: Mapping
metadata:
  name: petstore
  namespace: default
spec:
  prefix: "/"
  service: petstore.default:80
  rewrite: ""
  cors:
    origins: http://foo.example,http://bar.example
    methods: POST,GET,OPTIONS
    headers: Content-Type
    exposed_headers: X-Custom-Header,X-Other-Custom-Header
    credentials: false
    max_age: "120"
`,
	},
	{
		name: "cors-path-override",
		options: options.Options{
			Namespace: "default",
			Service: options.ServiceOptions{
				Namespace: "default",
				Name:      "petstore",
			},
			Ingress: options.IngressOptions{
				CORS: options.CORSOptions{
					Origins:       []string{"http://foo.example", "http://bar.example"},
					Methods:       []string{"POST", "GET", "OPTIONS"},
					Headers:       []string{"Content-Type"},
					ExposeHeaders: []string{"X-Custom-Header", "X-Other-Custom-Header"},
					Credentials:   nil,
					MaxAge:        120,
				},
			},
			PathSubOptions: map[string]options.SubOptions{
				"/pet": {
					CORS: options.CORSOptions{
						Origins:       []string{"http://bar.example"},
						Methods:       []string{"POST"},
						Headers:       []string{"Content-Type"},
						ExposeHeaders: []string{"X-Custom-Header", "X-Other-Custom-Header"},
						Credentials:   nil,
						MaxAge:        240,
					},
				},
			},
		},
		spec: `
openapi: 3.0.2
info:
  title: Swagger Petstore - OpenAPI 3.0
  version: 1.0.5
paths:
  "/pet":
    put:
      operationId: updatePet
      responses:
        '200':
          description: Successful operation
  "/pet/{petId}/uploadImage":
    post:
      operationId: uploadFile
      parameters:
        - name: petId
          in: path
          description: ID of pet to update
          required: true
          schema:
            type: integer
            format: int64
      responses:
        '200':
          description: Successful operation`,
		res: `
---
apiVersion: getambassador.io/v2
kind: Mapping
metadata:
  name: petstore-updatepet
  namespace: default
spec:
  prefix: "/pet"
  method: PUT
  service: petstore.default:80
  rewrite: ""
  cors:
    origins: http://bar.example
    methods: POST
    headers: Content-Type
    exposed_headers: X-Custom-Header,X-Other-Custom-Header
    credentials: false
    max_age: "240"
---
apiVersion: getambassador.io/v2
kind: Mapping
metadata:
  name: petstore-uploadfile
  namespace: default
spec:
  prefix: "/pet/([a-zA-Z0-9]*)/uploadImage"
  prefix_regex: true
  method: POST
  service: petstore.default:80
  rewrite: ""
  cors:
    origins: http://foo.example,http://bar.example
    methods: POST,GET,OPTIONS
    headers: Content-Type
    exposed_headers: X-Custom-Header,X-Other-Custom-Header
    credentials: false
    max_age: "120"
`,
	},
	{
		name: "timeouts-global",
		options: options.Options{
			Namespace: "default",
			Service: options.ServiceOptions{
				Namespace: "default",
				Name:      "petstore",
			},
			Timeouts: options.TimeoutOptions{
				RequestTimeout: 42,
				IdleTimeout:    43,
			},
		},
		spec: `
openapi: 3.0.2
info:
  title: Swagger Petstore - OpenAPI 3.0
  version: 1.0.5
paths:
  "/pet":
    put:
      operationId: updatePet
      responses:
        '200':
          description: Successful operation
  "/pet/{petId}/uploadImage":
    post:
      operationId: uploadFile
      parameters:
        - name: petId
          in: path
          description: ID of pet to update
          required: true
          schema:
            type: integer
            format: int64
      responses:
        '200':
          description: Successful operation`,
		res: `
---
apiVersion: getambassador.io/v2
kind: Mapping
metadata:
  name: petstore
  namespace: default
spec:
  prefix: "/"
  service: petstore.default:80
  rewrite: ""
  timeout_ms: 42000
  idle_timeout_ms 43000
`,
	},
	{
		name: "timeouts-path-override",
		options: options.Options{
			Namespace: "default",
			Service: options.ServiceOptions{
				Namespace: "default",
				Name:      "petstore",
			},
			Timeouts: options.TimeoutOptions{
				RequestTimeout: 42,
				IdleTimeout:    43,
			},
			PathSubOptions: map[string]options.SubOptions{
				"/pet": {
					Timeouts: options.TimeoutOptions{
						RequestTimeout: 35,
						IdleTimeout:    36,
					},
				},
			},
		},
		spec: `
openapi: 3.0.2
info:
  title: Swagger Petstore - OpenAPI 3.0
  version: 1.0.5
paths:
  "/pet":
    put:
      operationId: updatePet
      responses:
        '200':
          description: Successful operation
  "/pet/{petId}/uploadImage":
    post:
      operationId: uploadFile
      parameters:
        - name: petId
          in: path
          description: ID of pet to update
          required: true
          schema:
            type: integer
            format: int64
      responses:
        '200':
          description: Successful operation`,
		res: `
---
apiVersion: getambassador.io/v2
kind: Mapping
metadata:
  name: petstore-updatepet
  namespace: default
spec:
  prefix: "/pet"
  method: PUT
  service: petstore.default:80
  rewrite: ""
  timeout_ms: 35000
  idle_timeout_ms 36000
---
apiVersion: getambassador.io/v2
kind: Mapping
metadata:
  name: petstore-uploadfile
  namespace: default
spec:
  prefix: "/pet/([a-zA-Z0-9]*)/uploadImage"
  prefix_regex: true
  method: POST
  service: petstore.default:80
  rewrite: ""
  timeout_ms: 42000
  idle_timeout_ms 43000
`,
	},
}

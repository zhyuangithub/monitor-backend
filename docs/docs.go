// Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/ChangeNodeName": {
            "post": {
                "description": "修改某个节点的名称",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Post"
                ],
                "summary": "修改节点名称",
                "parameters": [
                    {
                        "description": "JSON数据",
                        "name": "Data",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/uiserver.changeNodeNameStruct"
                        }
                    }
                ],
                "responses": {}
            }
        },
        "/GetNodeDataHistory": {
            "get": {
                "description": "获取节点的数据记录",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Get"
                ],
                "summary": "获取数据记录",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "起始时间戳",
                        "name": "startTimestamp",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "结束时间戳",
                        "name": "endTimestamp",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "节点id",
                        "name": "nodeId",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {}
            }
        },
        "/GetNodeEventLogs": {
            "get": {
                "description": "获取单个节点的事件信息",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Get"
                ],
                "summary": "获取节点事件",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "节点id",
                        "name": "nodeId",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "事件类型分类 -1:all; 0:dismount; 1:nfc; 2:sleeping",
                        "name": "category",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "数量",
                        "name": "count",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "页码",
                        "name": "page",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {}
            }
        },
        "/GetNodeInfo": {
            "get": {
                "description": "获取单个节点的信息，节点状态字符串种类：online,offline",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Get"
                ],
                "summary": "获取节点信息",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "节点id",
                        "name": "nodeId",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {}
            }
        },
        "/GetNodesInfo": {
            "get": {
                "description": "按页码和数量获取全部节点的信息，节点状态字符串种类：online,offline",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Get"
                ],
                "summary": "批量获取节点信息",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "数量",
                        "name": "count",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "页码",
                        "name": "page",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "节点名称，可进行模糊查询",
                        "name": "name",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "节点状态，online,offline两种",
                        "name": "status",
                        "in": "query"
                    }
                ],
                "responses": {}
            }
        },
        "/notificationCenter": {
            "get": {
                "description": "websocket通知接口",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Websocket"
                ],
                "summary": "websocket通知",
                "responses": {}
            }
        }
    },
    "definitions": {
        "uiserver.changeNodeNameStruct": {
            "type": "object",
            "properties": {
                "nodeId": {
                    "type": "integer"
                },
                "nodeName": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
# Anytype API Documentation

## API Overview

- **Base URL**: `http://localhost:31009/v1`
- **Current API Version**: `2025-03-17` (passed in HTTP header `Anytype-Version`)
- **Authentication**: Bearer token (obtained through pairing process)
- **Default Timeout**: 10 seconds
- **Content Type**: All requests must use `application/json`

### Standard Headers

All requests should include:
```
Authorization: Bearer YOUR_TOKEN
Anytype-Version: 2025-03-17
Content-Type: application/json
```

### Error Handling

Errors follow a standard format:
```json
{
  "message": "Error description",
  "error": "Error code"
}
```

Common error scenarios:
- Invalid/expired authentication token
- Missing required parameters
- Invalid request format
- Version compatibility issues
- Resource not found
- Permission denied

### Pagination

All collection endpoints support pagination with consistent parameters:
- `offset`: Starting position (zero-based)
- `limit`: Number of items per page
- Response includes total count and has_more flag
- Default limits vary by endpoint (commonly 20-100 items)
- Maximum limit may be enforced by server

Example pagination URL:
```
/v1/spaces/{spaceId}/objects?offset=20&limit=10
```

### Common Response Formats

Standard response structure for collections:
```json
{
  "items": [
    {
      "id": "string",
      "name": "string",
      "type": "string"
    }
  ],
  "total": 0,
  "limit": 0,
  "offset": 0
}
```

### Common Data Structures

#### Object Structure
```json
{
  "id": "string",
  "type": "string",
  "name": "string",
  "icon": "string",
  "snippet": "string",
  "layout": "string",
  "space_id": "string",
  "root_id": "string",
  "props": {
    "key": "value"
  },
  "relations": {
    "tags": ["string"],
    "other_relations": "value"
  }
}
```

#### Member Object
```json
{
  "id": "string",
  "type": "string",
  "name": "string",
  "icon": "string",
  "role": "string",
  "identity": "string",
  "global_name": "string"
}
```

## Authentication Flow

1. **Display Code**
   - Client requests a display code
   - Desktop app shows a 4-digit code
   - Code valid for limited time (5 minutes)

2. **Verify Code**
   - Client submits the code shown in desktop app
   - If code matches, access token is issued
   - Token valid for extended period (30 days)

3. **Use Token**
   - Include token in all subsequent requests
   - Refresh or re-authenticate if token expires

## API Endpoints

### Authentication

- `POST /v1/auth/display_code` - Request a display code
- `POST /v1/auth/token` - Exchange code for token

### Spaces

- `GET /v1/spaces` - List spaces
- `GET /v1/spaces/{spaceId}` - Get space details
- `POST /v1/spaces` - Create space

### Objects

- `GET /v1/spaces/{spaceId}/objects` - List objects in space
- `GET /v1/spaces/{spaceId}/objects/{objectId}` - Get object details
- `POST /v1/spaces/{spaceId}/objects` - Create object
- `PATCH /v1/spaces/{spaceId}/objects/{objectId}` - Update object
- `DELETE /v1/spaces/{spaceId}/objects/{objectId}` - Delete object

### Types

- `GET /v1/spaces/{spaceId}/types` - List types
- `GET /v1/spaces/{spaceId}/types/{typeId}` - Get type details

### Templates

- `GET /v1/spaces/{spaceId}/types/{typeId}/templates` - List templates for type
- `GET /v1/spaces/{spaceId}/types/{typeId}/templates/{templateId}` - Get template

### Lists

- `GET /v1/spaces/{spaceId}/lists/{listId}/views` - Get list views
- `GET /v1/spaces/{spaceId}/lists/{listId}/view/{viewId}/objects` - Get objects in list
- `POST /v1/spaces/{spaceId}/lists/{listId}/objects` - Add objects to list
- `DELETE /v1/spaces/{spaceId}/lists/{listId}/objects/{objectId}` - Remove object from list

### Members

- `GET /v1/spaces/{spaceId}/members` - List members
- `GET /v1/spaces/{spaceId}/members/{memberId}` - Get member details
- `PATCH /v1/spaces/{spaceId}/members/{memberId}` - Update member

### Search

- `POST /v1/search` - Global search
- `POST /v1/spaces/{spaceId}/search` - Search within space

### Export

- `GET /v1/spaces/{spaceId}/objects/{objectId}/pdf` - Export as PDF
- `GET /v1/spaces/{spaceId}/objects/{objectId}/markdown` - Export as Markdown
- `GET /v1/spaces/{spaceId}/objects/{objectId}/html` - Export as HTML

## Search Parameters

The search endpoints support advanced filtering:

```json
{
  "query": "search text",
  "types": ["type1", "type2"],
  "tags": ["tag1", "tag2"],
  "filter": "custom filter expression",
  "sort": "sort expression",
  "limit": 100,
  "offset": 0,
  "custom": {}
}
```

Search features:
- Full text search across objects
- Type filtering
- Tag-based filtering
- Custom filter expressions
- Sorting options
- Default limit: 100 items
- Default offset: 0

## Advanced Features

### Type System

Types have the following structure:
```json
{
  "type": "string",
  "id": "string",
  "unique_key": "string",
  "name": "string",
  "icon": "string",
  "recommended_layout": "string",
  "relations": [
    {
      "key": "string",
      "name": "string",
      "format": "string"
    }
  ]
}
```

Type relationships:
- Each object must have a type
- Types define available relations
- Types can have recommended layouts
- Types support custom relations

### Template System

Templates provide structured starting points for objects:
```json
{
  "id": "string",
  "name": "string",
  "description": "string",
  "content": "string",
  "default_relations": {
    "key": "value"
  }
}
```

Template features:
- Support variable substitution (e.g., {{Title}}, {{Date}})
- Can include default relations
- May specify layout preferences
- Support markdown formatting
- Can be type-specific

### Object Relations

Relations system:
- Objects can have multiple relations
- Relations are typed (tags, links, etc.)
- Relations can be bidirectional
- Support complex queries and filtering
- Can be used for organizing and categorizing

Example relation query:
```json
{
  "relations": {
    "tags": {
      "$in": ["important", "work"]
    },
    "status": {
      "$eq": "active"
    }
  }
}
```

### Integration Guidelines

Best practices for API integration:
1. Connection Management
   - Maintain persistent connections
   - Handle connection errors gracefully
   - Implement exponential backoff
   - Monitor API availability

2. Data Synchronization
   - Track object modifications
   - Handle concurrent updates
   - Implement conflict resolution
   - Cache frequently accessed data

3. Performance Optimization
   - Use appropriate page sizes
   - Implement request batching
   - Cache authentication tokens
   - Monitor response times

4. Error Recovery
   - Implement retry logic
   - Handle token expiration
   - Log failed requests
   - Maintain audit trail

### API Version Management

Version handling:
- Version specified via `Anytype-Version` header
- Current version: `2025-03-17`
- Desktop app version compatibility checks
- Breaking changes require version update
- Multiple versions may be supported simultaneously

### HTTP Methods and Status Codes

The API uses standard HTTP methods:
- `GET`: Retrieve resources
- `POST`: Create resources or perform actions
- `PATCH`: Partial resource updates
- `DELETE`: Remove resources
- `OPTIONS`: Check allowed operations

### Security Considerations

- API only accessible from localhost (31009)
- Authentication required for all endpoints
- Token-based authorization
- No cross-origin requests allowed
- Rate limiting may be applied
- Desktop app must be running

### Testing and Development

Testing recommendations:
- Use debug mode for detailed logs
- Include version header in all requests
- Check error responses
- Handle pagination properly
- Validate object types
- Test with various permission levels

Example debug configuration:
```json
{
  "debug": true,
  "timeout": 30,
  "logLevel": "debug"
}
```

# Anytype API Examples with curl

Below are examples of how to interact with each Anytype API endpoint using curl commands. For all examples, replace:
- `YOUR_TOKEN` with your actual bearer token
- `SPACE_ID`, `OBJECT_ID`, etc. with actual IDs

## Authentication

### Display Code
```bash
curl -X POST "http://localhost:31009/v1/auth/display_code?app_name=RaycastExtension" \
  -H "Anytype-Version: 2025-03-17"
```
Response:
```json
{
  "challenge_id": "abc123xyz"
}
```

### Get Token
```bash
curl -X POST "http://localhost:31009/v1/auth/token?challenge_id=abc123xyz&code=1234" \
  -H "Anytype-Version: 2025-03-17"
```
Response:
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

## Spaces

### Get Spaces
```bash
curl "http://localhost:31009/v1/spaces?offset=0&limit=10" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Anytype-Version: 2025-03-17"
```
Response:
```json
{
  "items": [
    {
      "id": "space123",
      "name": "Personal Space",
      "icon": "📓"
    }
  ],
  "total": 1,
  "limit": 10,
  "offset": 0
}
```

### Get Space
```bash
curl "http://localhost:31009/v1/spaces/space123" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Anytype-Version: 2025-03-17"
```
Response:
```json
{
  "id": "space123",
  "name": "Personal Space",
  "icon": "📓",
  "created_at": "2025-01-15T12:30:00Z",
  "updated_at": "2025-03-10T16:45:22Z"
}
```

### Create Space
```bash
curl -X POST "http://localhost:31009/v1/spaces" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Anytype-Version: 2025-03-17" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Project X",
    "icon": "🚀"
  }'
```
Response:
```json
{
  "id": "space456",
  "name": "Project X",
  "icon": "🚀",
  "created_at": "2025-04-11T14:22:10Z"
}
```

## Objects

### Get Objects
```bash
curl "http://localhost:31009/v1/spaces/space123/objects?offset=0&limit=20" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Anytype-Version: 2025-03-17"
```
Response:
```json
{
  "items": [
    {
      "id": "obj1",
      "title": "My First Note",
      "type_id": "note"
    },
    {
      "id": "obj2",
      "title": "Project Roadmap",
      "type_id": "document"
    }
  ],
  "total": 42,
  "limit": 20,
  "offset": 0
}
```

### Get Object
```bash
curl "http://localhost:31009/v1/spaces/space123/objects/obj1" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Anytype-Version: 2025-03-17"
```
Response:
```json
{
  "id": "obj1",
  "title": "My First Note",
  "type_id": "note",
  "created_at": "2025-02-10T08:15:30Z",
  "updated_at": "2025-03-05T16:20:45Z",
  "content": "This is the content of my first note...",
  "relations": {
    "tags": ["important", "personal"]
  }
}
```

### Create Object
```bash
curl -X POST "http://localhost:31009/v1/spaces/space123/objects" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Anytype-Version: 2025-03-17" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Meeting Notes",
    "type_id": "note",
    "content": "Topics to discuss: 1. Project timeline 2. Resource allocation",
    "relations": {
      "tags": ["meeting", "important"]
    }
  }'
```
Response:
```json
{
  "id": "obj3",
  "title": "Meeting Notes",
  "type_id": "note"
}
```

### Delete Object
```bash
curl -X DELETE "http://localhost:31009/v1/spaces/space123/objects/obj3" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Anytype-Version: 2025-03-17"
```
Response:
```json
{
  "success": true
}
```

## Types

### Get Types
```bash
curl "http://localhost:31009/v1/spaces/space123/types?offset=0&limit=10" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Anytype-Version: 2025-03-17"
```
Response:
```json
{
  "items": [
    {
      "id": "note",
      "name": "Note",
      "icon": "📝"
    },
    {
      "id": "task",
      "name": "Task",
      "icon": "✅"
    }
  ],
  "total": 8,
  "limit": 10,
  "offset": 0
}
```

### Get Type
```bash
curl "http://localhost:31009/v1/spaces/space123/types/note" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Anytype-Version: 2025-03-17"
```
Response:
```json
{
  "id": "note",
  "name": "Note",
  "icon": "📝",
  "relations": [
    {
      "key": "tags",
      "name": "Tags",
      "format": "tag"
    }
  ]
}
```

## Templates

### Get Templates
```bash
curl "http://localhost:31009/v1/spaces/space123/types/note/templates?offset=0&limit=5" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Anytype-Version: 2025-03-17"
```
Response:
```json
{
  "items": [
    {
      "id": "template1",
      "name": "Basic Note",
      "description": "Simple note template"
    },
    {
      "id": "template2",
      "name": "Meeting Notes",
      "description": "Template for meeting notes"
    }
  ],
  "total": 2,
  "limit": 5,
  "offset": 0
}
```

### Get Template
```bash
curl "http://localhost:31009/v1/spaces/space123/types/note/templates/template1" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Anytype-Version: 2025-03-17"
```
Response:
```json
{
  "id": "template1",
  "name": "Basic Note",
  "description": "Simple note template",
  "content": "# {{Title}}\n\n_Created on {{Date}}_\n\n## Notes\n\n",
  "default_relations": {
    "tags": []
  }
}
```

## Lists

### Get List Views
```bash
curl "http://localhost:31009/v1/spaces/space123/lists/list1/views?offset=0&limit=5" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Anytype-Version: 2025-03-17"
```
Response:
```json
{
  "items": [
    {
      "id": "view1",
      "name": "Default View",
      "type": "grid"
    },
    {
      "id": "view2",
      "name": "Timeline",
      "type": "calendar"
    }
  ],
  "total": 2,
  "limit": 5,
  "offset": 0
}
```

### Get Objects in List
```bash
curl "http://localhost:31009/v1/spaces/space123/lists/list1/view1/objects?offset=0&limit=10" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Anytype-Version: 2025-03-17"
```
Response:
```json
{
  "items": [
    {
      "id": "obj1",
      "title": "My First Note",
      "type_id": "note"
    },
    {
      "id": "obj2",
      "title": "Project Roadmap",
      "type_id": "document"
    }
  ],
  "total": 15,
  "limit": 10,
  "offset": 0
}
```

### Add Objects to List
```bash
curl -X POST "http://localhost:31009/v1/spaces/space123/lists/list1/objects" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Anytype-Version: 2025-03-17" \
  -H "Content-Type: application/json" \
  -d '{
    "object_ids": ["obj4", "obj5"]
  }'
```
Response:
```json
{
  "success": true,
  "added": 2
}
```

### Remove Objects from List
```bash
curl -X DELETE "http://localhost:31009/v1/spaces/space123/lists/list1/objects/obj4" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Anytype-Version: 2025-03-17"
```
Response:
```json
{
  "success": true
}
```

## Members

### Get Members
```bash
curl "http://localhost:31009/v1/spaces/space123/members?offset=0&limit=10" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Anytype-Version: 2025-03-17"
```
Response:
```json
{
  "items": [
    {
      "id": "user1",
      "name": "Alice",
      "role": "admin"
    },
    {
      "id": "user2",
      "name": "Bob",
      "role": "editor"
    }
  ],
  "total": 2,
  "limit": 10,
  "offset": 0
}
```

### Get Member
```bash
curl "http://localhost:31009/v1/spaces/space123/members/user1" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Anytype-Version: 2025-03-17"
```
Response:
```json
{
  "id": "user1",
  "name": "Alice",
  "role": "admin",
  "joined_at": "2025-01-15T12:30:00Z"
}
```

### Update Member
```bash
curl -X PATCH "http://localhost:31009/v1/spaces/space123/members/user2" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Anytype-Version: 2025-03-17" \
  -H "Content-Type: application/json" \
  -d '{
    "role": "admin"
  }'
```
Response:
```json
{
  "id": "user2",
  "name": "Bob",
  "role": "admin"
}
```

## Search

### Global Search
```bash
curl -X POST "http://localhost:31009/v1/search?offset=0&limit=20" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Anytype-Version: 2025-03-17" \
  -H "Content-Type: application/json" \
  -d '{
    "query": "project",
    "types": ["note", "document"]
  }'
```
Response:
```json
{
  "items": [
    {
      "id": "obj2",
      "title": "Project Roadmap",
      "type_id": "document",
      "space_id": "space123"
    },
    {
      "id": "obj6",
      "title": "Project Timeline",
      "type_id": "note",
      "space_id": "space456"
    }
  ],
  "total": 2,
  "limit": 20,
  "offset": 0
}
```

### Space Search
```bash
curl -X POST "http://localhost:31009/v1/spaces/space123/search?offset=0&limit=10" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Anytype-Version: 2025-03-17" \
  -H "Content-Type: application/json" \
  -d '{
    "query": "meeting",
    "types": ["note"]
  }'
```
Response:
```json
{
  "items": [
    {
      "id": "obj3",
      "title": "Meeting Notes",
      "type_id": "note"
    }
  ],
  "total": 1,
  "limit": 10,
  "offset": 0
}
```

## Export

### Export Object
```bash
curl "http://localhost:31009/v1/spaces/space123/objects/obj1/pdf" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Anytype-Version: 2025-03-17" \
  -o "exported-note.pdf"
```

This will download the object as a PDF file named "exported-note.pdf" to your current directory.

Alternate formats may include:
```bash
# Export as markdown
curl "http://localhost:31009/v1/spaces/space123/objects/obj1/markdown" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Anytype-Version: 2025-03-17" \
  -o "exported-note.md"

# Export as HTML
curl "http://localhost:31009/v1/spaces/space123/objects/obj1/html" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Anytype-Version: 2025-03-17" \
  -o "exported-note.html"
```

# API Reference - DataFlow Platform

## Overview

The DataFlow Platform REST API provides programmatic access to manage data pipelines, transformations, and analytics workflows. This document covers all available endpoints, request/response formats, and usage examples.

## Base URL

```
https://api.dataflow-platform.io/v1
```

## Authentication

All API requests require authentication using Bearer tokens. Include your API key in the Authorization header:

```
Authorization: Bearer YOUR_API_KEY
```

## Endpoints

### Pipelines

#### List All Pipelines

```http
GET /pipelines
```

**Query Parameters:**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| page | integer | No | Page number for pagination (default: 1) |
| limit | integer | No | Results per page (default: 20, max: 100) |
| status | string | No | Filter by status: `active`, `paused`, `failed` |
| sort | string | No | Sort field: `created_at`, `updated_at`, `name` |

**Response:**

```json
{
  "data": [
    {
      "id": "pipe_1a2b3c4d",
      "name": "Customer Data ETL",
      "status": "active",
      "created_at": "2025-01-15T10:30:00Z",
      "updated_at": "2025-01-23T14:22:00Z",
      "source": {
        "type": "postgresql",
        "connection_id": "conn_xyz789"
      },
      "destination": {
        "type": "snowflake",
        "connection_id": "conn_abc123"
      }
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 20,
    "total": 45,
    "pages": 3
  }
}
```

#### Create Pipeline

```http
POST /pipelines
```

**Request Body:**

```json
{
  "name": "New Data Pipeline",
  "source_id": "conn_source123",
  "destination_id": "conn_dest456",
  "schedule": "0 */6 * * *",
  "transformations": [
    {
      "type": "filter",
      "config": {
        "column": "status",
        "operator": "equals",
        "value": "active"
      }
    }
  ]
}
```

**Response:** `201 Created`

### Connections

#### Get Connection Details

```http
GET /connections/:connection_id
```

**Response:**

```json
{
  "id": "conn_abc123",
  "name": "Production Database",
  "type": "postgresql",
  "host": "db.example.com",
  "port": 5432,
  "database": "production",
  "status": "connected",
  "last_tested": "2025-01-24T09:15:00Z"
}
```

## Error Handling

The API uses conventional HTTP response codes:

| Status Code | Meaning |
|-------------|---------|
| 200 | Success |
| 201 | Created |
| 400 | Bad Request - Invalid parameters |
| 401 | Unauthorized - Invalid API key |
| 403 | Forbidden - Insufficient permissions |
| 404 | Not Found |
| 429 | Too Many Requests - Rate limit exceeded |
| 500 | Internal Server Error |

**Error Response Format:**

```json
{
  "error": {
    "code": "invalid_parameter",
    "message": "The 'schedule' field must be a valid cron expression",
    "field": "schedule"
  }
}
```

## Rate Limiting

API requests are limited to 1000 requests per hour per API key. Check response headers for limit information:

```
X-RateLimit-Limit: 1000
X-RateLimit-Remaining: 847
X-RateLimit-Reset: 1706096400
```

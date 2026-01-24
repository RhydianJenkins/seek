# Getting Started with DataFlow Platform

## Introduction

Welcome to DataFlow Platform! This guide will help you set up your account, create your first data pipeline, and understand the core concepts of our platform.

## What is DataFlow Platform?

DataFlow Platform is a cloud-native data integration and transformation service that enables you to:

- Connect to multiple data sources (databases, APIs, file storage)
- Transform and clean data using visual or code-based tools
- Schedule automated data syncs
- Monitor data quality and pipeline health
- Scale processing to handle millions of records

## Prerequisites

Before you begin, ensure you have:

- A DataFlow Platform account (sign up at https://dataflow-platform.io)
- Access credentials for your data source
- Basic understanding of data structures (tables, JSON, CSV)

## Step 1: Create Your Account

1. Visit https://dataflow-platform.io/signup
2. Enter your email and create a password
3. Verify your email address
4. Complete the onboarding questionnaire
5. Choose your pricing plan

## Step 2: Connect Your First Data Source

### Connecting a Database

1. Navigate to **Connections** in the left sidebar
2. Click **New Connection**
3. Select your database type (PostgreSQL, MySQL, MongoDB, etc.)
4. Enter connection details:
   - **Host**: Your database host address
   - **Port**: Database port (e.g., 5432 for PostgreSQL)
   - **Database Name**: Name of the database
   - **Username**: Database user
   - **Password**: Database password

5. Click **Test Connection** to verify
6. Save your connection

### Security Note

All connection credentials are encrypted at rest using AES-256 encryption. Connections use SSL/TLS by default.

## Step 3: Create Your First Pipeline

A pipeline defines how data flows from a source to a destination, with optional transformations in between.

### Using the Visual Builder

1. Click **Pipelines** → **New Pipeline**
2. Name your pipeline (e.g., "Customer Data Sync")
3. Select your source connection
4. Choose the table or data to sync
5. Select or create a destination
6. Configure sync schedule (real-time, hourly, daily, custom cron)
7. Click **Create Pipeline**

### Example: Daily Customer Sync

```
Source: PostgreSQL (customers table)
   ↓
Transform: Filter active customers only
   ↓
Transform: Rename columns to match schema
   ↓
Destination: Snowflake (analytics.customers)
   ↓
Schedule: Daily at 2:00 AM UTC
```

## Step 4: Add Transformations (Optional)

Transformations let you modify data as it flows through your pipeline:

- **Filter**: Remove rows based on conditions
- **Map**: Rename or reformat columns
- **Aggregate**: Group and summarize data
- **Join**: Combine data from multiple sources
- **Custom**: Write JavaScript or Python code

### Example Filter Transformation

```javascript
// Keep only customers from the US
row => row.country === 'US'
```

## Step 5: Monitor Your Pipeline

Once your pipeline is running:

1. View the **Dashboard** for overall health
2. Click on your pipeline to see:
   - Run history
   - Records processed
   - Success/failure rate
   - Average execution time
3. Set up alerts for failures or data quality issues

## Next Steps

Now that you have your first pipeline running:

- Explore **Transformations** to clean and enrich data
- Set up **Webhooks** for event-driven workflows
- Configure **Data Quality Rules** to validate data
- Read the **API Reference** for programmatic access
- Join our community forum for tips and best practices

## Common Issues

### Connection Timeout

If your database connection times out, check:
- Firewall rules allow DataFlow IP addresses
- Database is accessible from external networks
- Credentials are correct

### Pipeline Fails on First Run

- Verify source and destination schemas match
- Check transformation logic for errors
- Review error logs in pipeline details

## Support

- Documentation: https://docs.dataflow-platform.io
- Community Forum: https://community.dataflow-platform.io
- Email Support: support@dataflow-platform.io
- Status Page: https://status.dataflow-platform.io

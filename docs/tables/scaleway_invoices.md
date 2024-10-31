---
title: "Steampipe Table: scaleway_invoice - Query Scaleway Invoices using SQL"
description: "Allows users to query Scaleway Invoices, providing detailed information about billing and usage for Scaleway services."
---

# Table: scaleway_invoice - Query Scaleway Invoices using SQL

Scaleway Invoices are detailed records of charges for the use of Scaleway's cloud services.

## Table Usage Guide

The `scaleway_invoice` table provides insights into billing information within Scaleway. As a finance manager or cloud administrator, explore invoice-specific details through this table, including total amounts, billing periods, and associated organizations. Utilize it to track expenses, verify charges, and manage cloud spending across different projects and timeframes.

## Examples

### Basic info
Explore the basic details of your Scaleway invoices, including their unique identifiers, associated organizations, and billing periods. This can help in tracking and managing your cloud expenses effectively.

```sql
SELECT
  id,
  organization_id,
  billing_period,
  total_taxed,
  state
FROM
  scaleway_invoice;
```

### List unpaid invoices
Identify any outstanding payments by listing all unpaid invoices. This query helps in managing financial obligations and ensuring timely payments.

```sql
SELECT
  id,
  organization_id,
  billing_period,
  total_taxed,
  due_date
FROM
  scaleway_invoice
WHERE
  state = 'unpaid';
```

### Get total billed amount for each organization
Calculate the total amount billed to each organization. This provides an overview of cloud spending across different entities within your Scaleway account.

```sql
SELECT
  organization_id,
  SUM(CAST(total_taxed AS DECIMAL)) as total_billed
FROM
  scaleway_invoice
GROUP BY
  organization_id;
```

### Find invoices with high discount amounts
Identify invoices with significant discounts. This can help in understanding which billing periods or services are providing the most cost savings.

```sql
SELECT
  id,
  billing_period,
  total_discount,
  total_taxed
FROM
  scaleway_invoice
WHERE
  CAST(total_discount AS DECIMAL) > 100
ORDER BY
  CAST(total_discount AS DECIMAL) DESC;
```

### List invoices for a specific date range
Retrieve invoices within a specific time frame. This is useful for periodic financial reviews or audits.

```sql
SELECT
  id,
  billing_period,
  total_taxed,
  issued_date
FROM
  scaleway_invoice
WHERE
  issued_date BETWEEN '2023-01-01' AND '2023-12-31'
ORDER BY
  issued_date;
```

### Get the average invoice amount by month
Calculate the average invoice amount for each month. This helps in understanding monthly spending patterns and budgeting for cloud services.

```sql
SELECT
  DATE_TRUNC('month', issued_date) AS month,
  AVG(CAST(total_taxed AS DECIMAL)) AS average_invoice_amount
FROM
  scaleway_invoice
GROUP BY
  DATE_TRUNC('month', issued_date)
ORDER BY
  month;
```

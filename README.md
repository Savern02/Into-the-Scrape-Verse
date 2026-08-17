# Into-the-Scrape-Verse
You write a scraper, it works, and a week later the site changes its layout and everything breaks quietly. Build one that repairs itself instead, run it from your coding agent, and spend the week turning the data into something real. https://www.wemakedevs.org/hackathons/scrape-verse/resourcesv

Every submission must include:

    A public source-code repository
    A clear README
    Example structured output
    A demo video showing the working project
    A clear explanation of how Bright Data Scraper Studio is used


Tech Stack:
Frontend: Next.JS + React + TailwindCss
Backend: Go, PostgreSQL 

| Component      | What it does                | Where it runs                  |
| -------------- | --------------------------- | -----------------------------  |
| Next.js        | Web app + server-side logic | **Vercel**                     |
| Go             | Data processing/workers     | **VPS**                        |
| PostgreSQL     | Canonical structured data   | **Managed PostgreSQL**         |
| Object storage | Raw Bright Data data        | **Managed object storage, S3** |

but for the hackathon:
| Component   | Free option                   | What runs there                       |
| ----------- | ----------------------------- | ------------------------------------- |
| Web app     | **Vercel**                    | Next.js frontend + server             |
| PostgreSQL  | **Supabase Free**             | Canonical data                        |
| Raw storage | **Cloudflare R2**             | Bright Data raw files                 |
| Go pipeline | **Google Cloud free-tier VM** | Go workers/scheduler                  |
| Scraping    | Bright Data                   | Your existing scraping infrastructure |



Current architecture: 
Price Comparison Platform — Architecture
Overview

The application is built around a data-processing pipeline rather than treating the project as a traditional web application.

The core architecture is:

Retailer Websites
       │
       ▼
 Bright Data
       │
       ▼
 Object Storage
   (Raw Data)
       │
       ▼
 Go Data Pipeline
 ├── Validation
 ├── Filtering
 ├── Normalization
 ├── Product Matching
 └── Derived Calculations
       │
       ▼
 PostgreSQL
 (Canonical Data)
       │
       ▼
    Next.js
   Web App
       │
       ▼
     Users
Components
Bright Data

Responsible for collecting data from retailer websites.

Bright Data provides the raw source data. The application does not treat scraped data as automatically correct.

Object Storage

Stores the raw Bright Data responses.

Raw data is preserved so that the processing pipeline can be changed or improved without requiring the data to be scraped again.

raw/
├── retailer-a/
│   └── 2026-08-17/
├── retailer-b/
│   └── 2026-08-17/
└── retailer-c/
    └── 2026-08-17/
Go Data Pipeline

The Go service is responsible for processing and maintaining the application's data.

Responsibilities include:

Data ingestion
Filtering
Validation
Normalization
Deduplication
Product matching
Price calculations
Data quality checks
Updating PostgreSQL

The Go pipeline runs independently of user requests and can be triggered by scheduled jobs or, as the system grows, a job queue.

PostgreSQL

PostgreSQL stores the canonical, validated representation of the data used by the application.

Examples include:

Products
Retailers
Offers
Prices
Price history
Categories
Product relationships

PostgreSQL is the application's source of truth for canonical data.

Next.js

Next.js is responsible for the web application.

Responsibilities include:

Search
Product pages
Price comparisons
Filtering and sorting
User accounts
Preferences
SEO
Price history visualization

Next.js consumes the canonical data produced by the Go pipeline.

Separation of Responsibilities

The system intentionally separates data acquisition, processing, storage, and presentation.

Bright Data
    │
    │ Collects
    ▼
Object Storage
    │
    │ Raw data
    ▼
Go Pipeline
    │
    │ Processes
    ▼
PostgreSQL
    │
    │ Canonical data
    ▼
Next.js
    │
    │ Presents
    ▼
Users
Key Principle

Raw observations are preserved; everything else is derived.

This allows transformations and business rules to evolve without necessarily requiring the data to be collected again.

For example, if we later introduce:

price_per_oz
price_per_lb
price_per_100g

we can update the transformation pipeline and reprocess existing raw data.

Data Independence

The Go pipeline and Next.js application operate independently.

If the Go pipeline is temporarily unavailable, Next.js can continue serving the latest known-good data from PostgreSQL.

If Next.js is unavailable, the Go pipeline can continue collecting and processing data.

┌─────────────────┐          ┌─────────────────┐
│   Go Pipeline   │          │     Next.js     │
│                 │          │                 │
│ Process Data   │          │ Serve Users     │
└────────┬────────┘          └────────┬────────┘
         │                            │
         │         PostgreSQL         │
         └──────────────┬─────────────┘
                        │
                        ▼
                 Canonical Data
Initial Technology Stack
Component	Technology	Purpose
Data Acquisition	Bright Data	Collect retailer data
Raw Storage	Object Storage	Preserve raw scraped data
Data Processing	Go	Validate and transform data
Database	PostgreSQL	Store canonical data
Web Application	Next.js	Serve the user-facing application
Future Scaling

The initial architecture should remain simple.

As data volume increases, the system can introduce:

Job queues
Multiple Go workers
Dedicated APIs
Analytics databases
Data warehouse infrastructure
More sophisticated product-matching systems

These can be added without fundamentally changing the core architecture.

Bright Data
     ↓
Object Storage
     ↓
Go Workers
     ↓
PostgreSQL
     ↓
Next.js
     ↓
Users

The data pipeline is the core of the platform; Next.js is the interface through which users consume the resulting data.

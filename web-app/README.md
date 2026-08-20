# Cheap Chick — Web App

The frontend for **Cheap Chick**, a grocery comparison platform that helps users find the cheapest prices across popular supermarkets.

## Overview

The Cheap Chick web application provides the user-facing interface for searching groceries, comparing prices, and creating grocery lists.

It is built with **Next.js, TypeScript, Tailwind CSS, shadcn/ui, and Supabase**.

## Features

* Grocery price comparison across supported supermarkets
* Search for individual grocery items
* Compare prices between stores
* Grocery list support
* Supabase authentication
* Responsive UI
* Server and client components with Next.js App Router
* Reusable UI components built with shadcn/ui
* Integration with the Cheap Chick data pipeline

## Tech Stack

* **Framework:** Next.js
* **Language:** TypeScript
* **Styling:** Tailwind CSS
* **UI:** shadcn/ui
* **Authentication:** Supabase Auth
* **Database / Storage:** Supabase
* **Runtime:** Node.js

## Project Structure

```text
web-app/
├── app/              # Next.js App Router pages and routes
├── components/       # Reusable UI components
├── lib/              # Shared utilities and Supabase configuration
├── public/           # Static assets
├── ...
```

## Getting Started

### Prerequisites

* Node.js
* npm
* A Supabase project

### Installation

From the `web-app` directory:

```bash
npm install
```

### Environment Variables

Create a `.env.local` file in the `web-app` directory:

```env
NEXT_PUBLIC_SUPABASE_URL=...
NEXT_PUBLIC_SUPABASE_ANON_KEY=...
```

These values can be found in your Supabase project's API settings.

### Run the Development Server

```bash
npm run dev
```

The application will be available at:

```text
http://localhost:3000
```

### Build for Production

```bash
npm run build
```

### Start the Production Server

```bash
npm start
```

## Authentication

Cheap Chick uses **Supabase Auth** for user authentication and session management.

Authentication is integrated with the Next.js App Router and allows authenticated users to access features such as grocery lists and other user-specific functionality.

## Data

The web application consumes grocery data produced by the Cheap Chick data pipeline.

The overall system is structured as:

```text
Supermarket Websites
        │
        ▼
   Site Scrapers
        │
        ▼
   Go Data Pipeline
        │
        ▼
 Data Validation / Ingest
        │
        ▼
      Supabase
        │
        ▼
   Cheap Chick Web App
```

## Development

The web application is designed to be modular so that new pages, components, grocery features, and supermarket integrations can be added without changing the core application structure.

When developing locally, make sure the required Supabase environment variables are configured before starting the application.

## Related Projects

The `web-app` is part of the larger Cheap Chick project.

```text
Cheap Chick
├── web-app/       # Next.js frontend
└── go-service/    # Data validation and ingestion pipeline
```

See the root repository README for the complete project architecture, scraping system, and data pipeline documentation.

## Acknowledgements

Built for the **WeMakeDevs Scrape-Verse Hackathon**.

## License

TBD — add your chosen license (e.g., MIT) or check the repository's `LICENSE` file.


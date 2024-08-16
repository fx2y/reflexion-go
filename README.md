# Reflexion Agent

## Overview

Reflexion Agent is a powerful tool-calling agent implemented using the Temporal Go SDK and instructor-ai/instructor-go. It demonstrates an advanced AI-powered workflow that combines initial response generation, self-reflection, and information retrieval to provide comprehensive and accurate answers to user queries.

## Features

- Implements the Reflexion architecture for improved AI responses
- Utilizes Temporal workflows for robust and scalable execution
- Integrates instructor-ai for AI-powered decision making
- Includes a mock search service for demonstration purposes
- Implements OpenTelemetry tracing for observability

## Implementation Details

The Reflexion Agent is implemented as a Temporal workflow with three main activities:

1. Initial Responder: Generates an initial answer, reflection, and search query.
2. Search: Performs a search based on the generated query.
3. Revisor: Revises the initial answer based on search results.

The workflow orchestrates these activities to produce a final, refined answer to the user's question.

## Configuration

- Update the `instructor.NewClient("your-api-key")` in `main.go` with your actual API key.
- Modify the `MockSearchService` in `mock_search.go` to integrate with a real search service if needed.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the LICENSE file for details.

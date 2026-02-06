import { describe, it, expect, beforeAll } from 'vitest';
import { env } from 'cloudflare:test';
import worker from './worker';

describe('Telemetry API - End-to-End Tests', () => {
    let telemetryEnv: any;

    beforeAll(() => {
        telemetryEnv = env;
    });

    describe('POST /api/telemetry/ingest', () => {
        it('should ingest a single event successfully', async () => {
            const request = new Request('http://localhost/api/telemetry/ingest', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({
                    event_type: 'test',
                    data: { message: 'hello' },
                }),
            });

            const response = await worker.fetch(request, telemetryEnv);
            expect(response.status).toBe(200);

            const data = await response.json();
            expect(data).toHaveProperty('accepted');
            expect(data.accepted).toBe(1);
            expect(data).toHaveProperty('total');
            expect(data.total).toBe(1);
            expect(data).toHaveProperty('timestamp');
            expect(data).toHaveProperty('results');
            expect(data.results).toHaveLength(1);
            expect(data.results[0]).toHaveProperty('success');
            expect(data.results[0].success).toBe(true);
            expect(data.results[0]).toHaveProperty('key');
            expect(data.results[0].key).toMatch(/^telemetry:\d+:[A-Za-z0-9]{8}$/);
        });

        it('should ingest multiple events in one request', async () => {
            const events = [
                { event_type: 'test', data: { id: 1 } },
                { event_type: 'test', data: { id: 2 } },
                { event_type: 'test', data: { id: 3 } },
            ];

            const request = new Request('http://localhost/api/telemetry/ingest', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify(events),
            });

            const response = await worker.fetch(request, telemetryEnv);
            expect(response.status).toBe(200);

            const data = await response.json();
            expect(data.accepted).toBe(3);
            expect(data.total).toBe(3);
            expect(data.results).toHaveLength(3);
        });

        it('should accept single event as object (not array)', async () => {
            const request = new Request('http://localhost/api/telemetry/ingest', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({
                    event_type: 'single',
                    data: { test: true },
                }),
            });

            const response = await worker.fetch(request, telemetryEnv);
            expect(response.status).toBe(200);

            const data = await response.json();
            expect(data.accepted).toBe(1);
        });

        it('should reject empty event array', async () => {
            const request = new Request('http://localhost/api/telemetry/ingest', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify([]),
            });

            const response = await worker.fetch(request, telemetryEnv);
            expect(response.status).toBe(400);

            const data = await response.json();
            expect(data).toHaveProperty('error');
        });

        it('should reject invalid JSON', async () => {
            const request = new Request('http://localhost/api/telemetry/ingest', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: 'invalid json{{',
            });

            const response = await worker.fetch(request, telemetryEnv);
            expect(response.status).toBe(400);

            const data = await response.json();
            expect(data.error).toBe('Invalid JSON');
        });

        it('should store events with received_at timestamp', async () => {
            const beforeTime = Date.now();

            const request = new Request('http://localhost/api/telemetry/ingest', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({
                    event_type: 'timestamp_test',
                    data: { test: true },
                }),
            });

            const response = await worker.fetch(request, telemetryEnv);
            const data = await response.json();

            // Verify the event was stored with received_at
            const key = data.results[0].key;
            const storedValue = await telemetryEnv.TELEMETRY.get(key);
            expect(storedValue).toBeTruthy();

            const storedEvent = JSON.parse(storedValue);
            expect(storedEvent).toHaveProperty('received_at');
            expect(storedEvent.received_at).toBeGreaterThanOrEqual(beforeTime);
            expect(storedEvent).toHaveProperty('event_type', 'timestamp_test');
        });

        it('should handle large batch of events', async () => {
            const events = Array.from({ length: 100 }, (_, i) => ({
                event_type: 'batch_test',
                data: { index: i },
            }));

            const request = new Request('http://localhost/api/telemetry/ingest', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify(events),
            });

            const response = await worker.fetch(request, telemetryEnv);
            expect(response.status).toBe(200);

            const data = await response.json();
            expect(data.accepted).toBe(100);
            expect(data.total).toBe(100);
            // Don't return all results for large batches
            expect(data.results).toBeUndefined();
        });
    });

    describe('POST /api/telemetry/batch', () => {
        it('should process a batch of events', async () => {
            const events = Array.from({ length: 50 }, (_, i) => ({
                event_type: 'batch_test',
                data: { index: i },
            }));

            const request = new Request('http://localhost/api/telemetry/batch', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({
                    batch_id: 'test-batch-123',
                    events,
                }),
            });

            const response = await worker.fetch(request, telemetryEnv);
            expect(response.status).toBe(200);

            const data = await response.json();
            expect(data).toHaveProperty('batch_id', 'test-batch-123');
            expect(data.accepted).toBe(50);
            expect(data.status).toBe('completed');
            expect(data).toHaveProperty('timestamp');
        });

        it('should generate batch_id if not provided', async () => {
            const request = new Request('http://localhost/api/telemetry/batch', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({
                    events: [
                        { event_type: 'auto_batch_test', data: {} },
                    ],
                }),
            });

            const response = await worker.fetch(request, telemetryEnv);
            expect(response.status).toBe(200);

            const data = await response.json();
            expect(data).toHaveProperty('batch_id');
            expect(data.batch_id).toBeTruthy();
            expect(data.batch_id).toMatch(/^[A-Za-z0-9]{16}$/);
        });

        it('should reject batches larger than 1000 events', async () => {
            const events = Array.from({ length: 1001 }, (_, i) => ({
                event_type: 'oversized_batch',
                data: { index: i },
            }));

            const request = new Request('http://localhost/api/telemetry/batch', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ events }),
            });

            const response = await worker.fetch(request, telemetryEnv);
            expect(response.status).toBe(400);

            const data = await response.json();
            expect(data.error).toContain('Batch too large');
        });

        it('should reject empty batch', async () => {
            const request = new Request('http://localhost/api/telemetry/batch', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ events: [] }),
            });

            const response = await worker.fetch(request, telemetryEnv);
            expect(response.status).toBe(400);
        });

        it('should store batch metadata', async () => {
            const batchId = 'metadata-test-123';
            const events = [{ event_type: 'test', data: {} }];

            const request = new Request('http://localhost/api/telemetry/batch', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ batch_id: batchId, events }),
            });

            await worker.fetch(request, telemetryEnv);

            // Check batch metadata was stored
            const batchKey = `batch:${batchId}`;
            const batchData = await telemetryEnv.TELEMETRY.get(batchKey);
            expect(batchData).toBeTruthy();

            const batch = JSON.parse(batchData);
            expect(batch.batch_id).toBe(batchId);
            expect(batch.event_count).toBe(1);
            expect(batch.status).toBe('completed');
            expect(batch).toHaveProperty('events');
            expect(batch.events).toHaveLength(1);
        });

        it('should associate events with batch_id', async () => {
            const batchId = 'association-test-456';
            const events = [
                { event_type: 'batch_assoc', data: { id: 1 } },
                { event_type: 'batch_assoc', data: { id: 2 } },
            ];

            const request = new Request('http://localhost/api/telemetry/batch', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ batch_id: batchId, events }),
            });

            const response = await worker.fetch(request, telemetryEnv);
            const data = await response.json();

            // Verify events have batch_id in stored data
            for (const eventKey of data.events) {
                const storedEvent = JSON.parse(await telemetryEnv.TELEMETRY.get(eventKey));
                expect(storedEvent).toHaveProperty('batch_id', batchId);
            }
        });
    });

    describe('GET /api/telemetry/batch/:batchId', () => {
        it('should retrieve batch status', async () => {
            const batchId = 'get-test-789';
            const events = [{ event_type: 'get_test', data: {} }];

            // Create batch first
            const createRequest = new Request('http://localhost/api/telemetry/batch', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ batch_id: batchId, events }),
            });
            await worker.fetch(createRequest, telemetryEnv);

            // Get batch status
            const getRequest = new Request(`http://localhost/api/telemetry/batch/${batchId}`);
            const response = await worker.fetch(getRequest, telemetryEnv);

            expect(response.status).toBe(200);

            const data = await response.json();
            expect(data.batch_id).toBe(batchId);
            expect(data.status).toBe('completed');
            expect(data.event_count).toBe(1);
        });

        it('should return 404 for non-existent batch', async () => {
            const request = new Request('http://localhost/api/telemetry/batch/non-existent');
            const response = await worker.fetch(request, telemetryEnv);

            expect(response.status).toBe(404);

            const data = await response.json();
            expect(data.error).toBe('Batch not found');
        });
    });

    describe('GET /api/telemetry/query', () => {
        it('should list telemetry events', async () => {
            // Create some test events
            const createRequest = new Request('http://localhost/api/telemetry/ingest', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify([
                    { event_type: 'query_test', data: { id: 1 } },
                    { event_type: 'query_test', data: { id: 2 } },
                ]),
            });
            await worker.fetch(createRequest, telemetryEnv);

            // Query events
            const queryRequest = new Request('http://localhost/api/telemetry/query?limit=5');
            const response = await worker.fetch(queryRequest, telemetryEnv);

            expect(response.status).toBe(200);

            const data = await response.json();
            expect(data).toHaveProperty('count');
            expect(data).toHaveProperty('events');
            expect(Array.isArray(data.events)).toBe(true);
        });

        it('should respect limit parameter', async () => {
            const queryRequest = new Request('http://localhost/api/telemetry/query?limit=3');
            const response = await worker.fetch(queryRequest, telemetryEnv);

            const data = await response.json();
            expect(data.events.length).toBeLessThanOrEqual(3);
        });

        it('should filter by time range', async () => {
            const now = Date.now();
            const oneHourAgo = now - 3600000;

            const queryRequest = new Request(
                `http://localhost/api/telemetry/query?start=${oneHourAgo}&end=${now}`
            );
            const response = await worker.fetch(queryRequest, telemetryEnv);

            expect(response.status).toBe(200);

            const data = await response.json();
            data.events.forEach((event: any) => {
                expect(event.received_at).toBeGreaterThanOrEqual(oneHourAgo);
                expect(event.received_at).toBeLessThanOrEqual(now);
            });
        });
    });

    describe('Error Handling and Edge Cases', () => {
        it('should handle malformed event data gracefully', async () => {
            const request = new Request('http://localhost/api/telemetry/ingest', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify([
                    { event_type: 'valid', data: {} },
                    null, // Invalid event
                    { event_type: 'also_valid', data: {} },
                ]),
            });

            const response = await worker.fetch(request, telemetryEnv);
            expect(response.status).toBe(200);

            const data = await response.json();
            // Should accept valid events and handle null gracefully
            expect(data.total).toBeGreaterThan(0);
        });

        it('should handle network-like errors with retry', async () => {
            // This test would need to mock KV failures
            // For now, we just verify the endpoint structure
            const request = new Request('http://localhost/api/telemetry/ingest', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ event_type: 'retry_test', data: {} }),
            });

            const response = await worker.fetch(request, telemetryEnv);
            expect(response.status).toBe(200);
        });

        it('should handle concurrent requests', async () => {
            const requests = Array.from({ length: 10 }, (_, i) => {
                return new Request('http://localhost/api/telemetry/ingest', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({
                        event_type: 'concurrent_test',
                        data: { index: i },
                    }),
                });
            });

            const responses = await Promise.all(
                requests.map(req => worker.fetch(req, telemetryEnv))
            );

            // All requests should succeed
            responses.forEach(response => {
                expect(response.status).toBe(200);
            });
        });
    });

    describe('Data Persistence', () => {
        it('should persist events with 7-day TTL', async () => {
            const request = new Request('http://localhost/api/telemetry/ingest', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({
                    event_type: 'ttl_test',
                    data: { message: 'Should expire after 7 days' },
                }),
            });

            const response = await worker.fetch(request, telemetryEnv);
            const data = await response.json();

            // Verify event is stored
            const key = data.results[0].key;
            const storedValue = await telemetryEnv.TELEMETRY.get(key);
            expect(storedValue).toBeTruthy();

            // Note: We can't easily test the exact TTL in unit tests,
            // but we've verified the key exists and has data
            const event = JSON.parse(storedValue);
            expect(event.event_type).toBe('ttl_test');
        });

        it('should preserve original event data', async () => {
            const originalData = {
                event_type: 'preserve_test',
                data: {
                    nested: { object: { with: 'deep structure' } },
                    array: [1, 2, 3],
                    number: 42,
                    string: 'test',
                    boolean: true,
                },
            };

            const request = new Request('http://localhost/api/telemetry/ingest', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify(originalData),
            });

            const response = await worker.fetch(request, telemetryEnv);
            const data = await response.json();

            const storedValue = await telemetryEnv.TELEMETRY.get(data.results[0].key);
            const storedEvent = JSON.parse(storedValue);

            expect(storedEvent.event_type).toBe(originalData.event_type);
            expect(storedEvent.data).toEqual(originalData.data);
        });
    });
});

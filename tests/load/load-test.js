import http from 'k6/http';
import { check, sleep } from 'k6';
import { Rate } from 'k6/metrics';

// Custom metrics
const errorRate = new Rate('errors');

// Test configuration
export const options = {
    stages: [
        { duration: '2m', target: 10 }, // Ramp up to 10 users
        { duration: '5m', target: 10 }, // Stay at 10 users
        { duration: '2m', target: 20 }, // Ramp up to 20 users
        { duration: '5m', target: 20 }, // Stay at 20 users
        { duration: '2m', target: 0 },  // Ramp down to 0 users
    ],
    thresholds: {
        http_req_duration: ['p(95)<500'], // 95% of requests must complete below 500ms
        http_req_failed: ['rate<0.1'],    // Error rate must be below 10%
        errors: ['rate<0.1'],             // Custom error rate must be below 10%
    },
};

const BASE_URL = 'http://localhost:8080';

export default function () {
    // Test 1: Health check
    const healthResponse = http.get(`${BASE_URL}/health`);
    check(healthResponse, {
        'health check status is 200': (r) => r.status === 200,
        'health check response time < 100ms': (r) => r.timings.duration < 100,
    });
    errorRate.add(healthResponse.status !== 200);

    // Test 2: Send event
    const eventPayload = JSON.stringify({
        type: 'user_action',
        user_id: `user_${Math.floor(Math.random() * 1000)}`,
        data: {
            action: 'test_action',
            page: `/page_${Math.floor(Math.random() * 10)}`,
            timestamp: new Date().toISOString(),
        },
        source: 'load_test',
    });

    const eventResponse = http.post(`${BASE_URL}/api/v1/events`, eventPayload, {
        headers: { 'Content-Type': 'application/json' },
    });

    check(eventResponse, {
        'event creation status is 200': (r) => r.status === 200,
        'event creation response time < 1000ms': (r) => r.timings.duration < 1000,
        'event creation returns event_id': (r) => {
            try {
                const body = JSON.parse(r.body);
                return body.event_id && body.event_id.length > 0;
            } catch (e) {
                return false;
            }
        },
    });
    errorRate.add(eventResponse.status !== 200);

    // Test 3: Generate random events
    const generateResponse = http.post(`${BASE_URL}/api/v1/generate?count=3`);
    check(generateResponse, {
        'generate events status is 200': (r) => r.status === 200,
        'generate events response time < 2000ms': (r) => r.timings.duration < 2000,
        'generate events returns count': (r) => {
            try {
                const body = JSON.parse(r.body);
                return body.count === 3;
            } catch (e) {
                return false;
            }
        },
    });
    errorRate.add(generateResponse.status !== 200);

    // Test 4: Get stats
    const statsResponse = http.get(`${BASE_URL}/api/v1/stats`);
    check(statsResponse, {
        'stats status is 200': (r) => r.status === 200,
        'stats response time < 500ms': (r) => r.timings.duration < 500,
        'stats returns data': (r) => {
            try {
                const body = JSON.parse(r.body);
                return body.total_events !== undefined;
            } catch (e) {
                return false;
            }
        },
    });
    errorRate.add(statsResponse.status !== 200);

    // Test 5: Get metrics
    const metricsResponse = http.get(`${BASE_URL}/metrics`);
    check(metricsResponse, {
        'metrics status is 200': (r) => r.status === 200,
        'metrics response time < 200ms': (r) => r.timings.duration < 200,
        'metrics returns prometheus format': (r) => r.body.includes('http_requests_total'),
    });
    errorRate.add(metricsResponse.status !== 200);

    // Random sleep between 1-3 seconds
    sleep(Math.random() * 2 + 1);
}

export function handleSummary(data) {
    return {
        'load-test-results.json': JSON.stringify(data, null, 2),
        stdout: `
Load Test Results:
==================
Total Requests: ${data.metrics.http_reqs.values.count}
Failed Requests: ${data.metrics.http_req_failed.values.count}
Error Rate: ${(data.metrics.http_req_failed.values.rate * 100).toFixed(2)}%
Average Response Time: ${data.metrics.http_req_duration.values.avg.toFixed(2)}ms
95th Percentile: ${data.metrics.http_req_duration.values['p(95)'].toFixed(2)}ms
Max Response Time: ${data.metrics.http_req_duration.values.max.toFixed(2)}ms
    `,
    };
}

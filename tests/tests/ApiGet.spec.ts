import {expect, test} from "@playwright/test";

test.describe('API tests', () => {
    test.beforeEach(async ({page}) => {
        await page.goto('http://localhost:8091/api/v1');
    });

    test('GET currencies', async ({request}) => {
        const response = await request.get('/currencies?scope=EUR');
        expect(response.status()).toBe(200);
    });

    test('GET ALL currencies', async ({request}) => {
        const response = await request.get('/currencies');
        expect(response.status()).toBe(200);
    });

    test('GET health test', async ({request}) => {
        const response = await request.get('/health');
        expect(response.status()).toBe(200);
    });
});
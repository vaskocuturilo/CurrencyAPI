import {expect, test} from "@playwright/test";

const DESCRIPTION =
    'A lightweight browser add-on that instantly converts any currency value on a webpage to your preferred currency. Automatically detects amounts, applies live exchange rates, and updates conversions in place across Google Chrome, Firefox, and Edge.';
const OG_DESCRIPTION =
    'A lightweight browser add-on that instantly converts any currency value on a webpage to your preferred currency. For Google Chrome, Firefox, and Edge.';

test.describe('Document head', () => {
    test.beforeEach(async ({page}) => {
        await page.goto('/');
    });

    test('has the expected meta description', async ({page}) => {
        await expect(page.locator('meta[name="description"]')).toHaveAttribute(
            'content',
            DESCRIPTION,
        );
    });

    test('og:description is the short social summary', async ({page}) => {
        await expect(
            page.locator('meta[property="og:description"]'),
        ).toHaveAttribute('content', OG_DESCRIPTION);
    });
});

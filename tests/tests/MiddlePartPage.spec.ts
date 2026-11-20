import {expect, test} from "@playwright/test";
import {MiddlePartPage} from "./pages/MiddlePartPage";

const DESCRIPTION =
    'A lightweight browser add-on that instantly converts any currency value on a webpage to your preferred currency. Automatically detects amounts, applies live exchange rates, and updates conversions in place across Google Chrome, Firefox, and Edge.';
const OG_DESCRIPTION =
    'Every price, in your currency.';

test.describe('Document head', () => {
    test.beforeEach(async ({page}) => {
        const middlePartPage = new MiddlePartPage(page)
        await middlePartPage.goto()
    });

    test('has the expected meta description', async ({page}) => {
        const middlePartPage = new MiddlePartPage(page)
        await expect(middlePartPage.descriptionText).toHaveText(DESCRIPTION)
        await expect(middlePartPage.ogDescriptionText).toHaveText(OG_DESCRIPTION)
    });
});

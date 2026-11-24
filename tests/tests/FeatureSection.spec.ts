import {expect, test} from "@playwright/test";
import {FeaturesSection} from "./pages/FeaturesSection";


const automaticDetectionText = "Finds prices and currency amounts on any webpage — no clicks or manual highlighting needed."
const liveExchangeRatesText = "Exchange rates refresh continuously so your converted prices always stay accurate."
const convertsInPlaceText = "Amounts update right where they appear on the page, keeping layouts clean and readable."
const lightweightAndPrivateText = "Tiny browser footprint. Detection and conversion run locally so your browsing stays on your device."

test.describe('Badges section', () => {

    test.beforeEach(async ({page}) => {
        const featuresSection = new FeaturesSection(page);
        await featuresSection.goto()
    });

    test('Check Features Section', async ({page}) => {
        const featuresSection = new FeaturesSection(page);

        await expect(featuresSection.titleFeatureSection).toBeVisible()
        await expect(featuresSection.automaticDetection).toContainText(automaticDetectionText)
        await expect(featuresSection.liveExchangeRates).toContainText(liveExchangeRatesText)
        await expect(featuresSection.convertsInPlace).toContainText(convertsInPlaceText)
        await expect(featuresSection.lightweightAndPrivate).toContainText(lightweightAndPrivateText)
    });
});
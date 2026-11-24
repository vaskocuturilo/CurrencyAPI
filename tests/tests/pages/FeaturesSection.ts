import {Locator, Page} from "@playwright/test";



export class FeaturesSection {
    readonly page: Page;
    readonly featureSection: Locator;
    readonly titleFeatureSection: Locator;
    readonly automaticDetection: Locator;
    readonly liveExchangeRates: Locator;
    readonly convertsInPlace: Locator;
    readonly lightweightAndPrivate: Locator;

    constructor(page: Page) {
        this.page = page;
        this.featureSection = page.locator("#features")
        this.titleFeatureSection = page.getByRole('heading', {name: 'Built for effortless browsing'})
        this.automaticDetection = page.getByRole('heading', {name: 'Automatic detection'})
        this.liveExchangeRates = page.getByRole('heading', {name: 'Live exchange rates'})
        this.convertsInPlace = page.getByRole('heading', {name: 'Converts in place'})
        this.lightweightAndPrivate = page.getByRole('heading', {name: 'Lightweight & private'})
    }

    async goto() {
        await this.page.goto('/');
    }
}
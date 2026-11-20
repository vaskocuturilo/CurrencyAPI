import {expect, Locator, Page} from "@playwright/test";

export class HeaderPage {
    readonly page: Page;
    readonly getCurrencyConverterLink: Locator;
    readonly getFeaturesLink: Locator;
    readonly getHowItWorksLink: Locator;
    readonly getFAQLink: Locator;
    readonly getAddToBrowserLink: Locator;


    constructor(page: Page) {
        this.page = page;
        this.getCurrencyConverterLink = page.locator('a', {hasText: 'Currency Converter'});
        this.getFeaturesLink = page.getByRole('link', {name: 'Features'});
        this.getHowItWorksLink = page.getByRole('link', {name: 'How it works'});
        this.getFAQLink = page.getByRole('link', {name: 'FAQ'});
        this.getAddToBrowserLink = page.getByRole('link', {name: 'Add to browser'});
    }

    async goto() {
        await this.page.goto('/');
    }

    async checkCurrencyConverterButton() {
        await expect(this.getCurrencyConverterLink).toBeVisible();
        await this.getCurrencyConverterLink.first().click();
    }
}
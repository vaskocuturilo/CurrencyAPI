import {Locator, Page} from "@playwright/test";

export class MiddlePartPage {
    readonly page: Page;
    readonly descriptionText: Locator;
    readonly ogDescriptionText: Locator;

    constructor(page: Page) {
        this.page = page;
        this.descriptionText = page.locator('p');
        this.ogDescriptionText = page.locator('h1').locator('span');
    }

    async goto() {
        await this.page.goto('/');
    }
}
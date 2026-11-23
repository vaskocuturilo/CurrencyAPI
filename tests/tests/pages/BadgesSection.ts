import {Locator, Page} from "@playwright/test";

export class BadgesSection {
    readonly page: Page;
    readonly getBadge: Locator;


    constructor(page: Page) {
        this.page = page;
        this.getBadge = page.locator('.rounded-full').getByRole("link");
    }

    async goto() {
        await this.page.goto('/');
    }
}
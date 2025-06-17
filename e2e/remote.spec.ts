import { test, expect } from "@playwright/test";

test.beforeEach(async ({ page }) => {
  await page.goto("http://limascope-with-agent:8080/");
});

test("has right title", async ({ page }) => {
  await expect(page).toHaveTitle(/.* - Limascope/);
});

test("select running container", async ({ page }) => {
  await page.getByTestId("side-menu").getByRole("link", { name: "limascope" }).click();
  await expect(page).toHaveURL(/\/container/);
  await expect(page.getByText("Accepting connections")).toBeVisible();
});

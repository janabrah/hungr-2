import { test, expect } from "@playwright/test";

// Helper to login before each test
async function login(
  page: import("@playwright/test").Page,
  email = "test@example.com",
) {
  await page.goto("/");
  await page.getByPlaceholder("Enter your email").fill(email);
  await page.getByRole("button", { name: "Login" }).click();
  await expect(page.getByRole("heading", { name: "Hungr" })).toBeVisible();
}

test.describe("Browse Recipes", () => {
  test.beforeEach(async ({ page }) => {
    await login(page);
  });

  test("can navigate to browse page", async ({ page }) => {
    await page.getByRole("link", { name: "Browse Recipes" }).click();
    await expect(
      page.getByRole("heading", { name: "Browse Recipes" }),
    ).toBeVisible();
  });

  test("browse page loads without errors", async ({ page }) => {
    await page.getByRole("link", { name: "Browse Recipes" }).click();

    // Should not show error message
    await expect(page.locator(".error")).not.toBeVisible();

    // Should show the recipe selector (even if empty)
    await expect(page.locator("select.select")).toBeVisible();
  });

  test("browse page passes email parameter to API", async ({ page }) => {
    // Intercept API calls to verify correct parameter
    let apiCallMade = false;
    let usedEmail = false;

    await page.route("**/api/recipes**", async (route) => {
      apiCallMade = true;
      const url = route.request().url();
      usedEmail = url.includes("email=");

      // Return empty response
      await route.fulfill({
        status: 200,
        contentType: "application/json",
        body: JSON.stringify({ recipeData: [], fileData: [] }),
      });
    });

    await page.getByRole("link", { name: "Browse Recipes" }).click();

    // Wait for API call
    await page.waitForTimeout(500);

    expect(apiCallMade).toBe(true);
    expect(usedEmail).toBe(true);
  });

  test("can filter recipes by tag", async ({ page }) => {
    await page.getByRole("link", { name: "Browse Recipes" }).click();

    const filterInput = page.getByPlaceholder("Filter by tag");
    await expect(filterInput).toBeVisible();

    await filterInput.fill("dinner");
    // Filter should be applied (URL should update)
    await expect(page).toHaveURL(/tag=dinner/);
  });
});

test.describe("Upload Recipes", () => {
  test.beforeEach(async ({ page }) => {
    await login(page);
  });

  test("can navigate to upload page", async ({ page }) => {
    await page.getByRole("link", { name: "Upload Recipe" }).click();
    await expect(
      page.getByRole("heading", { name: "Upload a Recipe" }),
    ).toBeVisible();
  });

  test("upload form has required fields", async ({ page }) => {
    await page.getByRole("link", { name: "Upload Recipe" }).click();

    await expect(page.locator('input[type="file"]')).toBeVisible();
    await expect(page.getByPlaceholder("Recipe name")).toBeVisible();
    await expect(page.getByPlaceholder("Tags (comma separated)")).toBeVisible();
    await expect(page.getByRole("button", { name: "Upload" })).toBeVisible();
  });

  test("upload passes email parameter to API", async ({ page }) => {
    let usedEmail = false;

    await page.route("**/api/recipes**", async (route) => {
      if (route.request().method() === "POST") {
        const url = route.request().url();
        usedEmail = url.includes("email=");

        await route.fulfill({
          status: 200,
          contentType: "application/json",
          body: JSON.stringify({
            success: true,
            recipe: {
              uuid: "test-uuid",
              name: "Test",
              user_uuid: "user-uuid",
              tag_string: "",
              created_at: new Date().toISOString(),
            },
            tags: [],
          }),
        });
      } else {
        await route.continue();
      }
    });

    await page.getByRole("link", { name: "Upload Recipe" }).click();

    // Fill out the form
    await page.getByPlaceholder("Recipe name").fill("Test Recipe");
    await page.getByPlaceholder("Tags (comma separated)").fill("test, e2e");

    // Create a test file
    const fileInput = page.locator('input[type="file"]');
    await fileInput.setInputFiles({
      name: "test.jpg",
      mimeType: "image/jpeg",
      buffer: Buffer.from("fake image data"),
    });

    // Submit
    await page.getByRole("button", { name: "Upload" }).click();

    // Wait for API call
    await page.waitForTimeout(500);

    expect(usedEmail).toBe(true);
  });
});

import Link from "next/link";

export default function Home() {
  return (
    <div className="grid grid-rows-[20px_1fr_20px] items-center justify-items-center min-h-screen p-8 pb-20 gap-16 sm:p-20 font-[family-name:var(--font-geist-sans)]">
      <main className="flex flex-col gap-8 row-start-2 items-center sm:items-start">
        <div>Welcome to Hungr!</div>
        <div>Would you like to upload a recipe or retrieve a recipe?</div>

        <div className="flex gap-4 items-center flex-col sm:flex-row">
          <Link href="/upload_recipe" passHref>
            <div className="rounded-full border border-solid border-transparent transition-colors flex items-center justify-center bg-foreground text-background gap-2 hover:bg-[#383838] dark:hover:bg-[#ccc] text-sm sm:text-base h-10 sm:h-12 px-4 sm:px-5 cursor-pointer">
              Upload Recipe
            </div>
          </Link>
          <Link href="/show_recipe" passHref>
            <div className="rounded-full border border-solid border-transparent transition-colors flex items-center justify-center bg-foreground text-background gap-2 hover:bg-[#383838] dark:hover:bg-[#ccc] text-sm sm:text-base h-10 sm:h-12 px-4 sm:px-5 cursor-pointer">
              Read Recipe
            </div>
          </Link>
        </div>
      </main>
    </div>
  );
}

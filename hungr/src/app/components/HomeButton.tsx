import Link from "next/link";
import Image from "next/image";

const HomeButton = () => {
  return (
    <Link href="/" passHref>
      <div className="inline-block bg-gray-800 text-white py-1 px-2 rounded hover:bg-gray-700 focus:outline-none focus:ring-2 focus:ring-gray-500 text-sm cursor-pointer">
        <Image src="/icon.png" alt="Home" width={46} height={46} />
      </div>
    </Link>
  );
};

export default HomeButton;

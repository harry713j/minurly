import { Link } from "react-router";

export default function NotFound() {
  return (
    <div className="flex flex-col justify-center items-center space-y-4">
      <h1 className="text-5xl font-semibold text-slate-700">404</h1>
      <h3 className="text-3xl font-medium text-slate-600">Page not found</h3>
      <Link
        className="text-lg transition delay-200 text-blue-400 hover:underline"
        to={"/dashboard"}
      >
        Dashboard
      </Link>
    </div>
  );
}

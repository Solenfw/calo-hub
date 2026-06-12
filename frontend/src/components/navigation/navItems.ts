import {
  BookOpen,
  CheckSquare,
  Home,
  RefreshCw,
  Scan,
} from "lucide-react";

export const navigationItems = [
  { href: "/", label: "Home", mobileLabel: "Home", icon: Home },
  {
    href: "/catalog",
    label: "Online Catalog",
    mobileLabel: "Catalog",
    icon: BookOpen,
  },
  {
    href: "/convert",
    label: "Quick Convert",
    mobileLabel: "Convert",
    icon: RefreshCw,
  },
  {
    href: "/extraction",
    label: "Extraction",
    mobileLabel: "Extract",
    icon: Scan,
  },
  {
    href: "/eselector",
    label: "eSelector",
    mobileLabel: "Select",
    icon: CheckSquare,
  },
] as const;

export function isNavigationActive(pathname: string, href: string) {
  if (href === "/") {
    return pathname === "/" || pathname === "/dashboard";
  }

  return pathname === href || pathname.startsWith(`${href}/`);
}

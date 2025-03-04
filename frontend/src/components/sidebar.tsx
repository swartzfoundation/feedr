import {
  Bookmark,
  CalendarDays,
  ChevronDownIcon,
  ChevronUp,
  Star,
  User2,
} from "lucide-react";

import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarGroup,
  SidebarGroupContent,
  SidebarGroupLabel,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
} from "./ui/sidebar";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "./ui/dropdown-menu";
import {
  Accordion,
  AccordionContent,
  AccordionItem,
  AccordionTrigger,
} from "./ui/accordion";

const items = [
  {
    title: "Today",
    url: "#",
    icon: CalendarDays,
  },
  {
    title: "Read Later",
    url: "#",
    icon: Bookmark,
  },
  {
    title: "Favourites",
    url: "#",
    icon: Star,
  },
];

const feeds = [
  {
    title: "People",
    items: [
      {
        title: "Chamath",
        link: "",
      },
      {
        title: "Kapoji",
        link: "",
      },
    ],
  },
  {
    title: "Technology",
    items: [
      {
        title: "Hacker News",
        link: "",
      },
      {
        title: "TechCrunch",
        link: "",
      },
      {
        title: "Wired",
        link: "",
      },
    ],
  },
  {
    title: "Health",
    items: [
      {
        title: "Nature Medicine",
        link: "",
      },
    ],
  },
  {
    title: "Travel",
    items: [
      {
        title: "Skift",
        link: "",
      },
      {
        title: "Phocuswright",
        link: "",
      },
    ],
  },
];

export function AppSidebar() {
  return (
    <Sidebar>
      <Content />
      <Footer />
    </Sidebar>
  );
}

function FeedGroup(props: { groups: typeof feeds }) {
  const groups = props.groups;
  const firstValue = groups[0].title;
  return (
    <Accordion type="single" defaultValue={firstValue}>
      {groups.map((g) => (
        <AccordionItem value={g.title}>
          <AccordionTrigger>
            <span>
              {g.title} {g.items.length}
            </span>
            <ChevronDownIcon className="AccordionChevron" aria-hidden />
          </AccordionTrigger>
          <AccordionContent>
            {g.items.map((item) => (
              <li>{item.title}</li>
            ))}
          </AccordionContent>
        </AccordionItem>
      ))}
    </Accordion>
  );
}

function Content() {
  return (
    <SidebarContent>
      <SidebarGroup>
        <SidebarGroupContent>
          <SidebarMenu>
            {items.map((item) => (
              <SidebarMenuItem key={item.title}>
                <SidebarMenuButton asChild>
                  <a href={item.url}>
                    <item.icon />
                    <span>{item.title}</span>
                  </a>
                </SidebarMenuButton>
              </SidebarMenuItem>
            ))}
          </SidebarMenu>
        </SidebarGroupContent>
      </SidebarGroup>
      <SidebarGroup>
        <SidebarGroupLabel>FEEDS</SidebarGroupLabel>
        <SidebarGroupContent>
          <SidebarMenu>
            <FeedGroup groups={feeds} />
          </SidebarMenu>
        </SidebarGroupContent>
      </SidebarGroup>
    </SidebarContent>
  );
}

function Footer() {
  return (
    <SidebarFooter>
      <SidebarMenu>
        <SidebarMenuItem>
          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <SidebarMenuButton>
                <User2 /> Username
                <ChevronUp className="ml-auto" />
              </SidebarMenuButton>
            </DropdownMenuTrigger>
            <DropdownMenuContent
              side="top"
              className="w-[--radix-popper-anchor-width]"
            >
              <DropdownMenuItem>
                <span>Account</span>
              </DropdownMenuItem>
              <DropdownMenuItem>
                <span>Billing</span>
              </DropdownMenuItem>
              <DropdownMenuItem>
                <span>Sign out</span>
              </DropdownMenuItem>
            </DropdownMenuContent>
          </DropdownMenu>
        </SidebarMenuItem>
      </SidebarMenu>
    </SidebarFooter>
  );
}

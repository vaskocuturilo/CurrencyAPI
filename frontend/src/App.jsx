import { ArrowRightLeft, Globe, Radar, Zap, ShieldCheck, ChevronDown } from "lucide-react";

const browsers = [
  { name: "Chrome", href: "#install" },
  { name: "Firefox", href: "#install" },
  { name: "Edge", href: "#install" },
];

const features = [
  {
    icon: Radar,
    t: "Automatic detection",
    d: "Finds prices and currency amounts on any webpage — no clicks or manual highlighting needed.",
  },
  {
    icon: Zap,
    t: "Live exchange rates",
    d: "Exchange rates refresh continuously so your converted prices always stay accurate.",
  },
  {
    icon: ArrowRightLeft,
    t: "Converts in place",
    d: "Amounts update right where they appear on the page, keeping layouts clean and readable.",
  },
  {
    icon: ShieldCheck,
    t: "Lightweight & private",
    d: "Tiny browser footprint. Detection and conversion run locally so your browsing stays on your device.",
  },
];

const faqs = [
  {
    q: "Is it free?",
    a: "Yes, the extension is completely free to install and use across Google Chrome, Firefox, and Edge.",
  },
  {
    q: "Which currencies are supported?",
    a: "All major world currencies, including USD, EUR, GBP, JPY, PLN, CAD, AUD, CHF, and many more.",
  },
  {
    q: "How often are rates updated?",
    a: "Live exchange rates are refreshed automatically throughout the day to reflect current market values.",
  },
  {
    q: "Does it collect my data?",
    a: "No. Page scanning and currency conversion happen locally in your browser without tracking your browsing history.",
  },
  {
    q: "Can I turn it off on certain sites?",
    a: "Yes, you can easily pause or toggle conversion for any website with a single click in the extension menu.",
  },
];

function InstallButtons() {
  return (
    <div className="flex flex-wrap gap-3">
      {browsers.map((b, i) => (
        <a
          key={b.name}
          href={b.href}
          className={
            i === 0
              ? "rounded-full bg-primary px-6 py-3 font-semibold text-primary-foreground shadow-soft transition hover:opacity-90"
              : "rounded-full border border-border bg-card px-6 py-3 font-semibold text-foreground transition hover:bg-secondary"
          }
        >
          Add to {b.name}
        </a>
      ))}
    </div>
  );
}

export function App() {
  return (
    <div className="min-h-screen font-sans">
      <header className="mx-auto flex max-w-6xl items-center justify-between px-6 py-5">
        <a href="#" className="flex items-center gap-2.5 font-display text-lg font-bold">
          <span className="grid h-9 w-9 place-items-center rounded-xl bg-primary text-primary-foreground">
            <ArrowRightLeft size={18} />
          </span>
          Currency Converter
        </a>
        <nav className="hidden gap-8 text-sm font-medium text-muted-foreground md:flex">
          <a href="#features" className="transition hover:text-foreground">
            Features
          </a>
          <a href="#how" className="transition hover:text-foreground">
            How it works
          </a>
          <a href="#faq" className="transition hover:text-foreground">
            FAQ
          </a>
        </nav>
        <a
          href="#install"
          className="rounded-full bg-accent px-5 py-2 text-sm font-semibold text-accent-foreground transition hover:opacity-90"
        >
          Add to browser
        </a>
      </header>

      <main>
        <section className="mx-auto grid max-w-6xl items-center gap-12 px-6 py-16 md:grid-cols-2 md:py-24">
          <div>
            <p className="mb-4 inline-block rounded-full bg-secondary px-3.5 py-1 text-xs font-semibold uppercase tracking-wider text-secondary-foreground">
              Lightweight Browser Add-on
            </p>
            <h1 className="font-display text-5xl font-extrabold leading-[1.05] tracking-tight md:text-6xl">
              Every price, <span className="text-primary">in your currency.</span>
            </h1>
            <p className="mt-6 max-w-lg text-lg leading-relaxed text-muted-foreground">
              A lightweight browser add‑on that instantly converts any currency value on a
              webpage to your preferred currency. It automatically detects amounts, applies live
              exchange rates, and updates conversions in place—making price comparison and
              international browsing effortless across Google Chrome, Firefox, and Edge.
            </p>
            <div className="mt-8">
              <InstallButtons />
            </div>
          </div>

          <div className="rounded-3xl border border-border bg-card p-5 shadow-soft">
            <div className="mb-4 flex items-center justify-between">
              <div className="flex gap-1.5">
                <span className="h-3 w-3 rounded-full bg-slate-200" />
                <span className="h-3 w-3 rounded-full bg-slate-200" />
                <span className="h-3 w-3 rounded-full bg-slate-200" />
              </div>
              <span className="rounded-full bg-accent px-2.5 py-0.5 text-xs font-medium text-accent-foreground">
                Auto-converted in place
              </span>
            </div>
            <div className="rounded-2xl bg-muted p-6">
              <div className="mb-4 flex h-36 items-center justify-center rounded-xl bg-secondary text-muted-foreground">
                <div className="text-center">
                  <div className="mx-auto mb-2 grid h-10 w-10 place-items-center rounded-full bg-card shadow-xs">
                    <ArrowRightLeft size={18} className="text-primary" />
                  </div>
                  <span className="text-xs font-medium">store.example.com/headphones</span>
                </div>
              </div>
              <p className="font-display text-lg font-semibold">Wireless Noise-Cancelling Headphones</p>
              <div className="relative mt-2 h-10 font-display text-3xl font-bold">
                <span className="animate-swap-a absolute">$49.99</span>
                <span className="animate-swap-b absolute text-primary">€46.20</span>
              </div>
              <p className="mt-3 text-xs text-muted-foreground">Live rate · 1 USD = 0.924 EUR</p>
            </div>
          </div>
        </section>

        <section id="features" className="mx-auto max-w-6xl px-6 py-16">
          <h2 className="font-display text-3xl font-bold tracking-tight md:text-4xl">
            Built for effortless browsing
          </h2>
          <div className="mt-10 grid gap-5 sm:grid-cols-2 lg:grid-cols-4">
            {features.map(({ icon: Icon, t, d }) => (
              <div
                key={t}
                className="rounded-2xl border border-border bg-card p-6 transition hover:shadow-soft"
              >
                <div className="inline-flex rounded-xl bg-accent/60 p-2.5 text-primary">
                  <Icon size={22} />
                </div>
                <h3 className="mt-4 font-display text-lg font-semibold">{t}</h3>
                <p className="mt-2 text-sm leading-relaxed text-muted-foreground">{d}</p>
              </div>
            ))}
          </div>
        </section>

        <section id="how" className="bg-primary py-20 text-primary-foreground">
          <div className="mx-auto max-w-6xl px-6">
            <h2 className="font-display text-3xl font-bold tracking-tight md:text-4xl">
              Up and running in seconds
            </h2>
            <ol className="mt-10 grid gap-8 md:grid-cols-3">
              {[
                "Install the extension in Google Chrome, Firefox, or Edge.",
                "Pick your preferred local currency.",
                "Browse normally — prices convert automatically in place.",
              ].map((s, i) => (
                <li key={s} className="rounded-2xl bg-white/5 p-6 backdrop-blur-xs">
                  <span className="font-display text-5xl font-extrabold text-accent">
                    0{i + 1}
                  </span>
                  <p className="mt-4 text-lg leading-snug">{s}</p>
                </li>
              ))}
            </ol>
          </div>
        </section>

        <section className="mx-auto max-w-6xl px-6 py-16 text-center">
          <p className="text-sm font-semibold uppercase tracking-wider text-muted-foreground">
            Works seamlessly on
          </p>
          <div className="mt-5 flex flex-wrap justify-center gap-4">
            {browsers.map((b) => (
              <span
                key={b.name}
                className="flex items-center gap-2 rounded-full border border-border bg-card px-6 py-2.5 font-display font-semibold shadow-2xs"
              >
                <Globe size={16} className="text-primary" />
                {b.name}
              </span>
            ))}
          </div>
        </section>

        <section id="faq" className="mx-auto max-w-3xl px-6 py-16">
          <h2 className="font-display text-3xl font-bold tracking-tight md:text-4xl">
            Frequently asked questions
          </h2>
          <div className="mt-8 divide-y divide-border rounded-2xl border border-border bg-card">
            {faqs.map((f) => (
              <details key={f.q} className="group p-5">
                <summary className="flex cursor-pointer list-none items-center justify-between font-display font-semibold">
                  <span>{f.q}</span>
                  <ChevronDown
                    size={18}
                    className="text-muted-foreground transition-transform duration-200 group-open:rotate-180"
                  />
                </summary>
                <p className="mt-3 text-sm leading-relaxed text-muted-foreground">{f.a}</p>
              </details>
            ))}
          </div>
        </section>

        <section id="install" className="mx-auto max-w-6xl px-6 pb-20">
          <div className="rounded-3xl bg-secondary p-10 text-center md:p-16">
            <h2 className="font-display text-3xl font-bold tracking-tight md:text-4xl">
              Stop converting in your head.
            </h2>
            <p className="mt-3 text-muted-foreground">
              Free and lightweight for Google Chrome, Firefox, and Edge.
            </p>
            <div className="mt-8 flex justify-center">
              <InstallButtons />
            </div>
          </div>
        </section>
      </main>

      <footer className="border-t border-border py-8 text-center text-sm text-muted-foreground">
        © {new Date().getFullYear()} Real-time Currency Converter Extension
      </footer>
    </div>
  );
}

export default App;

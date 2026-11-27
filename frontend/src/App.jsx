import { useState } from "react";
import { ArrowRightLeft, ArrowDownUp, Globe, Radar, Zap, ShieldCheck, ChevronDown } from "lucide-react";

const browsers = [
  { name: "Chrome", href: "#install" },
  { name: "Firefox", href: "#install" },
  { name: "Edge", href: "#install" },
];

const currencies = [
  { symbol: "$", code: "USD" },
  { symbol: "€", code: "EUR" },
  { symbol: "£", code: "GBP" },
  { symbol: "¥", code: "JPY" },
  { symbol: "Fr", code: "CHF" },
  { symbol: "$", code: "CAD" },
  { symbol: "$", code: "AUD" },
  { symbol: "¥", code: "CNY" },
  { symbol: "₹", code: "INR" },
  { symbol: "R$", code: "BRL" },
  { symbol: "$", code: "MXN" },
  { symbol: "₩", code: "KRW" },
  { symbol: "kr", code: "SEK" },
  { symbol: "zł", code: "PLN" },
  { symbol: "₺", code: "TRY" },
  { symbol: "R", code: "ZAR" },
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

function CurrencyTicker() {
  const repeated = [...currencies, ...currencies];
  return (
    <div className="overflow-hidden border-y border-[var(--ink)]/10 py-4" aria-hidden="true">
      <div
        className="animate-ticker flex w-max gap-10 whitespace-nowrap text-sm"
        style={{ fontFamily: "var(--font-mono)" }}
      >
        {repeated.map((c, idx) => (
          <span key={`${c.code}-${idx}`} className="flex items-center gap-2">
            <span className="font-semibold text-[var(--brass)]">{c.symbol}</span>
            <span className="tracking-widest">{c.code}</span>
          </span>
        ))}
      </div>
    </div>
  );
}

function HeroPreviewCard() {
  const [swapped, setSwapped] = useState(false);

  return (
    <div
      className="relative rounded-[1.75rem] p-7 shadow-[0_30px_60px_-30px_rgba(20,33,28,0.5)]"
      style={{ background: "var(--ink)", color: "var(--paper)" }}
    >
      <div className="flex items-center justify-between">
        <span
          className="text-xs uppercase tracking-[0.2em] text-[var(--paper)]/50"
          style={{ fontFamily: "var(--font-mono)" }}
        >
          Preview
        </span>
        <span className="flex gap-1.5" aria-hidden="true">
          <span className="size-2 rounded-full bg-[var(--paper)]/20" />
          <span className="size-2 rounded-full bg-[var(--paper)]/20" />
          <span className="size-2 rounded-full bg-[var(--brass)]" />
        </span>
      </div>

      <div className="mt-5 rounded-2xl bg-[var(--paper)]/[0.06] px-5 py-4">
        <div className="text-xs text-[var(--paper)]/50">You send</div>
        <div className="mt-1 flex items-center justify-between gap-4">
          <span
            className="text-3xl font-medium tabular-nums"
            style={{ fontFamily: "var(--font-mono)" }}
          >
            100.00
          </span>
          <span
            className="rounded-full bg-[var(--paper)] px-3 py-1 text-sm font-semibold tracking-widest text-[var(--ink)]"
            style={{ fontFamily: "var(--font-mono)" }}
          >
            {swapped ? "EUR" : "USD"}
          </span>
        </div>
      </div>

      <div className="relative my-1 flex justify-center">
        <div className="absolute inset-x-0 top-1/2 h-px bg-[var(--paper)]/10" />
        <button
          type="button"
          aria-label="Swap currencies"
          onClick={() => setSwapped((prev) => !prev)}
          className="relative grid size-11 cursor-pointer place-items-center rounded-full bg-[var(--brass)] text-[var(--ink)] transition-transform duration-500 hover:scale-105 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-[var(--paper)]"
          style={{ transform: swapped ? "rotate(180deg)" : "rotate(0deg)" }}
        >
          <ArrowDownUp className="size-5" aria-hidden="true" />
        </button>
      </div>

      <div className="mt-5 rounded-2xl bg-[var(--paper)]/[0.06] px-5 py-4">
        <div className="text-xs text-[var(--paper)]/50">They receive</div>
        <div className="mt-1 flex items-center justify-between gap-4">
          <span
            className="text-3xl font-medium tabular-nums text-[var(--paper)]/40"
            style={{ fontFamily: "var(--font-mono)" }}
          >
            —
          </span>
          <span
            className="rounded-full bg-[var(--paper)] px-3 py-1 text-sm font-semibold tracking-widest text-[var(--ink)]"
            style={{ fontFamily: "var(--font-mono)" }}
          >
            {swapped ? "USD" : "EUR"}
          </span>
        </div>
      </div>
      
      <CurrencyTicker />
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

          <HeroPreviewCard />
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

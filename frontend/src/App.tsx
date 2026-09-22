import { Calculator } from './components/Calculator';

export function App() {
  return (
    <main className="mx-auto flex min-h-screen w-full max-w-xl flex-col justify-center px-4 py-8 sm:px-6">
      <section className="rounded-2xl border border-slate-200 bg-white p-4 shadow-sm sm:p-6">
        <header className="mb-5">
          <h1 className="text-xl font-semibold text-slate-900 sm:text-2xl">Calculator</h1>
          <p className="mt-1 text-sm text-slate-600">
            Operations are provided by the calculator service.
          </p>
        </header>

        <Calculator />
      </section>
    </main>
  );
}

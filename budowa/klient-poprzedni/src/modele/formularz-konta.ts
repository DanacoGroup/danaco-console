import { AccountKind, type Account } from '../../../shared/contract';
import { przycisk, utworzWierszOdpowiedzi } from './kontrolki-formularza';
import { elementyPol, utworzPolaKonta } from './pola-konta';
import { nazwaRodzaju, rodzajUzywaKataloguKonfiguracji } from './rodzaje-kont';
import type { StanKont } from './stan-kont';

/**
 * Formularz jednego konta obsługuje założenie nowego konta albo zmianę
 * istniejącego, w tym pole poświadczenia wprowadzane wyłącznie jako wejście
 * oraz katalog konfiguracji właściwy dla rodzaju konta.
 */
export interface FormularzKonta {
  /** Formularz osadzany w panelu kont. */
  element: HTMLElement;
  /** Ustawia formularz na wskazane konto; `null` otwiera konto nowe. */
  pokaz(konto: Account | null): void;
}

export function utworzFormularzKonta(stan: StanKont): FormularzKonta {
  const pola = utworzPolaKonta();
  const { nazwa, rodzaj, znakRodzaju, dostawca, zewnetrzny, model } = pola;
  const { adres, katalog, poswiadczenie, czynne, domyslne } = pola;

  const zapis = przycisk('Zapisz konto', 'dn-btn dn-btn--atrament');
  const nowe = przycisk('Nowe konto', 'dn-btn dn-btn--zarys');
  const oznacz = przycisk('Uczyń domyślnym', 'dn-btn dn-btn--zarys');
  const usun = przycisk('Usuń konto', 'dn-btn dn-btn--niebezpieczny');
  const odpowiedz = utworzWierszOdpowiedzi();

  const przyciski = document.createElement('div');
  przyciski.className = 'dm-formularz__przyciski';
  przyciski.append(zapis, nowe, oznacz, usun);

  const tytul = document.createElement('h3');
  tytul.className = 'dm-formularz__tytul';

  const element = document.createElement('section');
  element.className = 'dm-formularz';
  element.append(tytul, ...elementyPol(pola), przyciski, odpowiedz.element);

  /** Konto pokazywane w formularzu; `null` znaczy konto nowe. */
  let biezace: Account | null = null;

  /** Rodzaj wskazany w formularzu — z listy przy nowym, z konta przy zmianie. */
  const wybranyRodzaj = (): string =>
    biezace === null ? rodzaj.kontrolka.value : biezace.kind;

  function ubierz(): void {
    const konto = biezace;
    const nowyByt = konto === null;
    tytul.textContent = konto === null ? 'Nowe konto' : `Konto: ${konto.name}`;
    rodzaj.element.hidden = !nowyByt;
    znakRodzaju.hidden = nowyByt;
    znakRodzaju.textContent = nowyByt ? '' : `Rodzaj: ${nazwaRodzaju(wybranyRodzaj())}`;
    katalog.element.hidden = !rodzajUzywaKataloguKonfiguracji(wybranyRodzaj());
    domyslne.element.hidden = !nowyByt;
    oznacz.hidden = konto === null || konto.isDefault;
    usun.hidden = nowyByt;
    nowe.hidden = nowyByt;
    poswiadczenie.kontrolka.placeholder = nowyByt
      ? 'poświadczenie nowego konta'
      : 'wpisz, aby zapisać nowe; puste zostawia dotychczasowe';
  }

  function wypelnij(konto: Account | null): void {
    biezace = konto;
    nazwa.kontrolka.value = konto?.name ?? '';
    dostawca.kontrolka.value = konto?.provider ?? '';
    zewnetrzny.kontrolka.value = konto?.externalId ?? '';
    model.kontrolka.value = konto?.defaultModel ?? '';
    adres.kontrolka.value = konto?.baseUrl ?? '';
    katalog.kontrolka.value = konto?.configDir ?? '';
    poswiadczenie.kontrolka.value = '';
    czynne.kontrolka.checked = konto?.enabled ?? true;
    domyslne.kontrolka.checked = false;
    odpowiedz.wyczysc();
    ubierz();
  }

  rodzaj.kontrolka.addEventListener('change', ubierz);

  zapis.addEventListener('click', () => void zapiszKonto());
  nowe.addEventListener('click', () => stan.wybierz(null));

  oznacz.addEventListener('click', () => {
    if (biezace === null) return;
    const zadanie = { accountId: biezace.id, kind: biezace.kind };
    void stan.ustawDomyslne(zadanie).then((wynik) => {
      odpowiedz.pokaz(
        wynik.udany
          ? 'Konto jest teraz domyślne dla swojego rodzaju.'
          : `Rdzeń nie przyjął oznaczenia: ${tresc(wynik.blad?.message)}`,
        wynik.udany,
      );
    });
  });

  usun.addEventListener('click', () => {
    if (biezace === null) return;
    void stan.usun(biezace.id).then((wynik) => {
      odpowiedz.pokaz(
        wynik.udany
          ? 'Konto usunięte. Kanały powołujące się na nie utraciły powiązanie i zostały w rejestrze.'
          : `Rdzeń nie usunął konta: ${tresc(wynik.blad?.message)}`,
        wynik.udany,
      );
    });
  });

  /** Zapis rozgałęzia się na dwie komendy: założenie albo zmianę. */
  async function zapiszKonto(): Promise<void> {
    const wynik =
      biezace === null
        ? await stan.dodaj({
            name: nazwa.kontrolka.value.trim(),
            kind: rodzaj.kontrolka.value as AccountKind,
            provider: dostawca.kontrolka.value.trim(),
            ...opcjonalne(),
            ...tajne(),
            makeDefault: domyslne.kontrolka.checked,
            enabled: czynne.kontrolka.checked,
          })
        : await stan.zmien({
            accountId: biezace.id,
            name: nazwa.kontrolka.value.trim(),
            provider: dostawca.kontrolka.value.trim(),
            ...opcjonalne(),
            ...tajne(),
            enabled: czynne.kontrolka.checked,
          });

    odpowiedz.pokaz(
      wynik.udany
        ? 'Konto zapisane. Poświadczenie, jeżeli podane, zostało przyjęte i nie wraca żadną odpowiedzią.'
        : `Rdzeń nie przyjął zapisu: ${tresc(wynik.blad?.message)}`,
      wynik.udany,
    );
    if (wynik.udany) poswiadczenie.kontrolka.value = '';
  }

  /** Pola nieobowiązkowe; puste pomijamy, żeby nie nadpisywać zapisu pustką. */
  function opcjonalne(): PolaOpcjonalne {
    const wartosci: PolaOpcjonalne = {};
    const zewnetrznyByt = zewnetrzny.kontrolka.value.trim();
    const nazwaModelu = model.kontrolka.value.trim();
    const adresBazowy = adres.kontrolka.value.trim();
    const katalogCLI = katalog.kontrolka.value.trim();

    if (zewnetrznyByt !== '') wartosci.externalId = zewnetrznyByt;
    if (nazwaModelu !== '') wartosci.defaultModel = nazwaModelu;
    if (adresBazowy !== '') wartosci.baseUrl = adresBazowy;
    if (katalogCLI !== '' && rodzajUzywaKataloguKonfiguracji(wybranyRodzaj())) {
      wartosci.configDir = katalogCLI;
    }
    return wartosci;
  }

  /** Poświadczenie wchodzi do żądania wyłącznie wtedy, gdy je wpisano. */
  function tajne(): { credential?: string } {
    const wartosc = poswiadczenie.kontrolka.value;
    return wartosc === '' ? {} : { credential: wartosc };
  }

  wypelnij(null);

  return { element, pokaz: wypelnij };
}

/** Pola opcjonalne konta, których pominięcie w żądaniu zmiany zostawia odpowiadający zapis po stronie rdzenia bez zmiany wartości. */
interface PolaOpcjonalne {
  externalId?: string;
  defaultModel?: string;
  baseUrl?: string;
  configDir?: string;
}

function tresc(powod: string | undefined): string {
  return powod ?? 'brak treści błędu';
}

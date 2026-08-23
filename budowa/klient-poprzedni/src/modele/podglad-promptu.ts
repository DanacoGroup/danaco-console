import { IdentityMode, type IdentityLayerContent } from '../../../shared/contract';
import { nazwaOsi } from '../konfiguracja/zasiegi';
import { przycisk } from './kontrolki-formularza';
import type { StanTozsamosci } from './stan-tozsamosci';
import { nazwaTrybu, nazwaWarstwy, ostrzezenieTrybu } from './warstwy-tozsamosci';

/**
 * Podgląd nakładki obowiązującej — treść, która trafi do modelu.
 *
 * Podgląd pokazuje, co złożyło się w jeden prompt i czy prompt fabryczny został
 * zastąpiony. Treść bierze się z komendy `identity.effective.get`, a nie ze
 * sklejenia warstw w kliencie: własny porządek składania rozjechałby się
 * z rdzeniem przy pierwszej zmianie reguł.
 * Kategoria wymagana bez treści nie wstrzymuje uruchomienia, więc podgląd
 * wylicza ją imiennie.
 */
export interface PodgladPromptu {
  /** Podgląd osadzany w sekcji modeli. */
  element: HTMLElement;
  /** Nanosi nakładkę odczytaną w stanie. */
  odswiez(): void;
}

export function utworzPodgladPromptu(stan: StanTozsamosci): PodgladPromptu {
  const naglowek = document.createElement('div');
  naglowek.className = 'dm-podglad__naglowek';

  const trybZnak = document.createElement('span');
  trybZnak.className = 'dn-plakietka dm-podglad__tryb';

  const odcisk = document.createElement('code');
  odcisk.className = 'dm-podglad__odcisk';

  const kopiuj = przycisk('Kopiuj złożony prompt', 'dn-btn dn-btn--sm dn-btn--zarys');
  naglowek.append(trybZnak, odcisk, kopiuj);

  const skutek = document.createElement('p');
  skutek.className = 'dm-edytor__ostrzezenie';
  skutek.setAttribute('role', 'note');

  const braki = document.createElement('p');
  braki.className = 'dm-podglad__braki';

  const warstwy = document.createElement('ol');
  warstwy.className = 'dm-podglad__warstwy';

  const prompt = document.createElement('pre');
  prompt.className = 'dm-podglad__prompt';

  const element = document.createElement('section');
  element.className = 'dm-podglad';
  element.append(naglowek, skutek, braki, tytulSekcji('Warstwy nakładki'), warstwy, tytulSekcji('Złożony prompt systemowy'), prompt);

  kopiuj.addEventListener('click', () => {
    // Schowek bywa niedostępny — brak uprawnienia albo kontekst bez zaufania.
    // Niepowodzenie zostaje w dzienniku; podgląd działa dalej.
    void navigator.clipboard
      ?.writeText(prompt.textContent ?? '')
      .catch((blad: unknown) => console.warn('[modele] schowek odmówił zapisu', blad));
  });

  function odswiez(): void {
    const nakladka = stan.nakladka();

    if (nakladka === null) {
      trybZnak.textContent = 'nakładka nieodczytana';
      trybZnak.dataset.zastapienie = 'false';
      odcisk.textContent = '';
      skutek.textContent = KOMUNIKAT_BEZ_NAKLADKI;
      skutek.dataset.zastapienie = 'false';
      braki.hidden = true;
      warstwy.replaceChildren();
      prompt.textContent = '';
      return;
    }

    const zastapienie = nakladka.mode === IdentityMode.ZASTAP;
    trybZnak.textContent = nazwaTrybu(nakladka.mode);
    trybZnak.className = zastapienie
      ? 'dn-plakietka dn-plakietka--ostrzezenie dm-podglad__tryb'
      : 'dn-plakietka dn-plakietka--informacja dm-podglad__tryb';
    trybZnak.dataset.zastapienie = String(zastapienie);
    odcisk.textContent = `odcisk: ${nakladka.promptHash}`;
    skutek.textContent = ostrzezenieTrybu(nakladka.mode);
    skutek.dataset.zastapienie = String(zastapienie);

    const brakujace = nakladka.missingRequiredCategoryIds ?? [];
    braki.hidden = brakujace.length === 0;
    braki.textContent = `Kategorie wymagane bez treści: ${brakujace.join(', ')}. Brak treści nie wstrzymuje uruchomienia, lecz nakładka pojedzie bez nich.`;

    warstwy.replaceChildren(
      ...(nakladka.layers.length > 0
        ? uporzadkuj(nakladka.layers).map(pozycjaWarstwy)
        : [komunikatBezWarstw()]),
    );

    prompt.textContent = nakladka.prompt;
  }

  return { element, odswiez };
}

/** Porządek warstw bierzemy z rdzenia; warstwa bez kolejności trafia na koniec. */
function uporzadkuj(warstwy: readonly IdentityLayerContent[]): IdentityLayerContent[] {
  return [...warstwy].sort((pierwsza, druga) => kolejnosc(pierwsza) - kolejnosc(druga));
}

function kolejnosc(warstwa: IdentityLayerContent): number {
  return Number.isFinite(warstwa.order) ? warstwa.order : Number.MAX_SAFE_INTEGER;
}

/** Jedna warstwa: skąd pochodzi treść i ile jej jest. */
function pozycjaWarstwy(warstwa: IdentityLayerContent): HTMLElement {
  const element = document.createElement('li');
  element.className = 'dm-podglad__warstwa';

  const nazwa = document.createElement('span');
  nazwa.className = 'dm-podglad__warstwa-nazwa';
  nazwa.textContent = `${warstwa.name} · ${nazwaWarstwy(warstwa.layer)}`;

  const zrodlo = document.createElement('span');
  zrodlo.className = 'dm-podglad__warstwa-zrodlo';
  zrodlo.textContent = opisPochodzenia(warstwa);

  element.append(nazwa, zrodlo);
  return element;
}

/** Pochodzenie warstwy: kategoria, oś oraz jej byt, gdy oś go ma. */
function opisPochodzenia(warstwa: IdentityLayerContent): string {
  const bytOsi = warstwa.axisId ?? '';
  const os = bytOsi === '' ? nazwaOsi(warstwa.axis) : `${nazwaOsi(warstwa.axis)} ${bytOsi}`;
  return `${warstwa.categoryId} · oś: ${os} · znaków: ${warstwa.content.length}`;
}

function tytulSekcji(tresc: string): HTMLElement {
  const element = document.createElement('h4');
  element.className = 'dm-podglad__tytul';
  element.textContent = tresc;
  return element;
}

function komunikatBezWarstw(): HTMLElement {
  const element = document.createElement('li');
  element.className = 'dm-podglad__warstwa dm-podglad__warstwa--pusta';
  element.textContent =
    'Nakładka nie ma ani jednej warstwy. Model dostanie prompt pusty albo — przy trybie DOŁĄCZ — sam prompt fabryczny.';
  return element;
}

const KOMUNIKAT_BEZ_NAKLADKI =
  'Nakładka obowiązująca nie została jeszcze odczytana z rdzenia. Odczytaj katalog ponownie; podgląd wypełni się, gdy rdzeń odpowie.';

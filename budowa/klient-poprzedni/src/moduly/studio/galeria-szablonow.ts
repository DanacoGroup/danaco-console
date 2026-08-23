import type { StudioTemplate } from '../../../../shared/contract';

/**
 * Galeria szablonów — miniatury pism, z których zakłada się dokument.
 *
 * ── Skąd szablony ───────────────────────────────────────────────────────────
 * Z rdzenia: `studio.template.list` oddaje pismo, umowę, raport, notatkę
 * i ofertę wraz z polami do wypełnienia, a `studio.template.apply` zakłada
 * dokument z wypełnieniem. Obie komendy mają uchwyt i obie działały, tylko nie
 * miały czym się pokazać. Galeria nie dopisuje ani jednego szablonu od siebie.
 *
 * ── Czym jest miniatura ─────────────────────────────────────────────────────
 * Podglądem WYGLĄDU złożonym z tego, co szablon o sobie mówi: nazwa,
 * przeznaczenie, format i pola. Obrazka miniatury kontrakt nie niesie i rdzeń
 * obrazów nie generuje, więc miniatura jest wyrysem układu — kartka, na niej
 * nagłówek nazwy i kreski wierszy w liczbie pól. Rysunek udający gotowe pismo
 * pokazywałby wygląd, którego szablon nie obiecuje.
 */

/** Czynności galerii zlecane oknu. */
export interface CzynnosciGalerii {
  /** Zakłada dokument z szablonu wraz z wartościami jego pól. */
  naZalozenie(idSzablonu: string, wartosci: Record<string, string>, tytul: string): void;
  /**
   * Otwiera szablon w warsztacie szablonów — pola, zmiana, oddanie do pliku.
   *
   * Nieobowiązkowe: galeria działa bez warsztatu, tylko wtedy wykaz jest wyłącznie
   * do CZYTANIA i zakładania dokumentu. Gdy warsztat jest zamontowany, miniatura
   * dostaje przejście do niego — bo szablon zmienia się tam, a nie tutaj.
   */
  naWarsztat?(idSzablonu: string): void;
}

/** Galeria wraz z jej odświeżeniem. */
export interface GaleriaSzablonow {
  element: HTMLElement;
  /** Wstawia szablony oddane przez rdzeń. */
  ustawSzablony(szablony: readonly StudioTemplate[]): void;
  /** Otwiera albo zamyka galerię. */
  przestawWidocznosc(): void;
  widoczny(): boolean;
}

export function utworzGalerieSzablonow(czynnosci: CzynnosciGalerii): GaleriaSzablonow {
  let otwarta = false;
  let wszystkie: readonly StudioTemplate[] = [];

  const szukanie = document.createElement('input');
  szukanie.type = 'search';
  szukanie.className = 'dn-pole-kontrolka';
  szukanie.placeholder = 'szukaj szablonu — nazwa albo przeznaczenie';
  szukanie.setAttribute('aria-label', 'Szukanie szablonu');

  const kategorie = document.createElement('div');
  kategorie.className = 'ms-galeria__kategorie';

  const miniatury = document.createElement('div');
  miniatury.className = 'ms-galeria__miniatury';

  const zdanie = document.createElement('p');
  zdanie.className = 'dn-pole-opis';

  const element = document.createElement('section');
  element.className = 'ms-galeria';
  element.hidden = true;
  element.setAttribute('aria-label', 'Galeria szablonów dokumentu');
  element.append(szukanie, kategorie, zdanie, miniatury);

  /**
   * Kategoria szablonu — fabryczny albo własny.
   *
   * Innej kategorii kontrakt nie niesie: `StudioTemplate` ma `builtin`, a nie
   * pole kategorii. Wymyślanie działów („pisma", „umowy") na podstawie nazw
   * byłoby porządkiem zgadniętym.
   */
  let kategoria: 'wszystkie' | 'fabryczne' | 'wlasne' = 'wszystkie';

  for (const pozycja of [
    { kod: 'wszystkie', nazwa: 'Wszystkie' },
    { kod: 'fabryczne', nazwa: 'Fabryczne' },
    { kod: 'wlasne', nazwa: 'Zapisane przez Operatora' },
  ] as const) {
    const przycisk = document.createElement('button');
    przycisk.type = 'button';
    przycisk.className = 'dn-btn dn-btn--sm dn-btn--duch';
    przycisk.textContent = pozycja.nazwa;
    przycisk.dataset['kategoria'] = pozycja.kod;
    przycisk.addEventListener('click', () => {
      kategoria = pozycja.kod;
      przerysuj();
    });
    kategorie.append(przycisk);
  }

  szukanie.addEventListener('input', () => przerysuj());

  function widoczneSzablony(): readonly StudioTemplate[] {
    const fraza = szukanie.value.trim().toLowerCase();
    return wszystkie.filter((szablon) => {
      if (kategoria === 'fabryczne' && !szablon.builtin) return false;
      if (kategoria === 'wlasne' && szablon.builtin) return false;
      if (fraza === '') return true;
      const opis = `${szablon.name} ${szablon.description ?? ''}`.toLowerCase();
      return opis.includes(fraza);
    });
  }

  function przerysuj(): void {
    const widoczne = widoczneSzablony();
    zdanie.textContent =
      wszystkie.length === 0
        ? 'Rdzeń nie oddał ani jednego szablonu — wykaz jest pusty, a nie nieodczytany.'
        : `Szablonów w wykazie ${wszystkie.length}, widocznych ${widoczne.length}. ` +
          'Zakładanie idzie komendą studio.template.apply i zakłada NOWY dokument sesji. Szablon ' +
          'własny zakłada się, zmienia i oddaje do pliku w warsztacie szablonów — wykaz fabryczny ' +
          'przestał być jedynym, jaki Operator ma.';
    miniatury.replaceChildren(...widoczne.map((szablon) => utworzMiniature(szablon, czynnosci)));
  }

  return {
    element,

    ustawSzablony(szablony) {
      wszystkie = szablony;
      przerysuj();
    },

    przestawWidocznosc() {
      otwarta = !otwarta;
      element.hidden = !otwarta;
    },

    widoczny: () => otwarta,
  };
}

/** Buduje jedną miniaturę wraz z polami szablonu i przyciskiem założenia. */
function utworzMiniature(szablon: StudioTemplate, czynnosci: CzynnosciGalerii): HTMLElement {
  const pola = szablon.fields ?? [];

  const kartka = document.createElement('div');
  kartka.className = 'ms-galeria__kartka';
  kartka.setAttribute('aria-hidden', 'true');
  const nazwaNaKartce = document.createElement('div');
  nazwaNaKartce.className = 'ms-galeria__kartka-naglowek';
  nazwaNaKartce.textContent = szablon.name;
  kartka.append(nazwaNaKartce);
  // Wiersz miniatury na każde pole szablonu: podgląd mówi, ile miejsc jest do
  // wypełnienia, i nic ponad to.
  for (const pole of pola) {
    const kreska = document.createElement('div');
    kreska.className = 'ms-galeria__kartka-wiersz';
    kreska.dataset['obowiazkowe'] = pole.required ? 'tak' : 'nie';
    kartka.append(kreska);
  }

  const nazwa = document.createElement('p');
  nazwa.className = 'ms-galeria__nazwa';
  nazwa.textContent = szablon.name;

  const opis = document.createElement('p');
  opis.className = 'dn-pole-opis';
  opis.textContent =
    `${szablon.description ?? 'Szablon bez opisu w rdzeniu.'} Format: ${szablon.format} · ` +
    `${szablon.builtin ? 'fabryczny' : 'zapisany przez Operatora'} · pól: ${pola.length}`;

  const tytul = document.createElement('input');
  tytul.type = 'text';
  tytul.className = 'dn-pole-kontrolka';
  tytul.placeholder = 'tytuł zakładanego dokumentu';
  tytul.setAttribute('aria-label', `Tytuł dokumentu z szablonu ${szablon.name}`);

  const kontrolki = new Map<string, HTMLInputElement>();
  const formularz = document.createElement('div');
  formularz.className = 'ms-galeria__pola';
  for (const pole of pola) {
    const kontrolka = document.createElement('input');
    kontrolka.type = 'text';
    kontrolka.className = 'dn-pole-kontrolka';
    kontrolka.placeholder = pole.required ? `${pole.label} (obowiązkowe)` : pole.label;
    kontrolka.value = pole.defaultValue ?? '';
    kontrolka.setAttribute('aria-label', pole.label);
    kontrolki.set(pole.name, kontrolka);
    formularz.append(kontrolka);
  }

  const zaloz = document.createElement('button');
  zaloz.type = 'button';
  zaloz.className = 'dn-btn dn-btn--sm dn-btn--sygnal';
  zaloz.textContent = 'Załóż dokument';
  zaloz.dataset['szablon'] = szablon.id;
  zaloz.addEventListener('click', () => {
    const wartosci: Record<string, string> = {};
    for (const [nazwaPola, kontrolka] of kontrolki) wartosci[nazwaPola] = kontrolka.value;
    czynnosci.naZalozenie(szablon.id, wartosci, tytul.value.trim());
  });

  const miniatura = document.createElement('article');
  miniatura.className = 'ms-galeria__miniatura';
  miniatura.dataset['szablon'] = szablon.id;
  miniatura.append(kartka, nazwa, opis, tytul, formularz, zaloz);

  const doWarsztatu = czynnosci.naWarsztat;
  if (doWarsztatu !== undefined) {
    const prowadz = document.createElement('button');
    prowadz.type = 'button';
    prowadz.className = 'dn-btn dn-btn--sm dn-btn--zarys';
    prowadz.textContent = 'Prowadź w warsztacie';
    prowadz.dataset['warsztat'] = szablon.id;
    prowadz.title = szablon.builtin
      ? 'Szablon fabryczny: pola da się w warsztacie przeczytać i oddać go do pliku, ale usunąć ' +
        'się go nie da — rdzeń odmawia nazwanym powodem. Zapisz z niego szablon własny i zmieniaj ' +
        'tamten.'
      : 'Otwiera szablon w warsztacie: pola do wypełnienia, zmiana, oddanie do pliku i usunięcie.';
    prowadz.addEventListener('click', () => doWarsztatu(szablon.id));
    miniatura.append(prowadz);
  }
  return miniatura;
}

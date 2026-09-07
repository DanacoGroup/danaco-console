// Nakładka Konfiguracji wywoływana z okna Ustawień: dokumenty tożsamości
// w czterech sekcjach oraz nastawy sesji wraz z ich zasięgiem.
import { Command, ConfigScope, SessionConfigArea } from '../../../shared/contract.ts';
import type { IdentityCategory, IdentityLayerContent } from '../../../shared/contract.ts';
import type { Kanal } from '../protokol/kanal.ts';
import { wywolaj } from '../protokol/wywolanie.ts';
import { oglos } from './ogloszenie.ts';
import { zapytajWOknieModalnym } from './pytanie-modalne.ts';

const NAGLOWEK = 'Konfiguracja';

/* Z czterech sekcji nakładki tylko Konstytucja niesie pole do pisania; Profil
   i Ekspertyza są podglądem warstw, a ostatnia sekcja trzyma nastawy sesji. */
const SEKCJA_DOKUMENTU = 'kf-sek-konstytucja';
const KOD_KATEGORII = 'konstytucja';

export function zwiazKonfiguracjeUstawien(kanal: Kanal, korzen: Element): void {
  const sekcja = korzen.querySelector('#sekcja-konto');
  if (sekcja === null) return;
  sekcja.append(przyciskWejscia(sekcja.ownerDocument));

  sekcja.addEventListener('click', (zdarzenie) => {
    const cel = zdarzenie.target;
    if (!(cel instanceof Element)) return;
    if (cel.closest('#us-konfiguracja-otworz') === null) return;
    zdarzenie.stopPropagation();
    void otworzKonfiguracje(kanal, korzen);
  });
}

function przyciskWejscia(dokument: Document): HTMLElement {
  const pozycja = dokument.createElement('div');
  pozycja.className = 'us-pozycja';
  const glowa = dokument.createElement('div');
  glowa.className = 'us-pozycja-glowa';
  const etykieta = dokument.createElement('span');
  etykieta.className = 'us-etykieta';
  etykieta.textContent = 'Konfiguracja platformy';
  const prawa = dokument.createElement('span');
  prawa.className = 'us-prawa';
  const przycisk = dokument.createElement('button');
  przycisk.type = 'button';
  przycisk.className = 'dn-btn dn-btn--atrament dn-btn--sm';
  przycisk.id = 'us-konfiguracja-otworz';
  przycisk.textContent = 'Otwórz Konfigurację';
  prawa.appendChild(przycisk);
  glowa.append(etykieta, prawa);
  const zdanie = dokument.createElement('p');
  zdanie.className = 'us-opis';
  zdanie.textContent = 'Zasady pracy modeli i nastawy sesji w jednym miejscu. '
    + 'Konfiguracja otwiera się nad bieżącym widokiem i zamyka bez jego zmiany.';
  pozycja.append(glowa, zdanie);
  return pozycja;
}

/* Nakładka stoi w warstwie okna, nie w sekcji: zamknięcie zdejmuje ją w całości,
   więc powtórne wejście zawsze zastaje świeżo wypełnione pola. */
async function otworzKonfiguracje(kanal: Kanal, korzen: Element): Promise<void> {
  const dokument = korzen.ownerDocument;
  if (dokument.querySelector('.kf-nakladka') !== null) return;
  const szablon = dokument.getElementById('dn-tresc-konfiguracja');
  if (!(szablon instanceof HTMLTemplateElement)) {
    oglos(NAGLOWEK, 'To wydanie nie niesie okna Konfiguracji.', 'ostrzezenie');
    return;
  }
  const wynik = await wywolaj(kanal, Command.ConfigWindowOpen, {
    scope: ConfigScope.Application,
  });
  if (!wynik.udany) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił otwarcia Konfiguracji.',
      'ostrzezenie');
    return;
  }
  const nakladka = szablon.content.cloneNode(true) as DocumentFragment;
  const korzenNakladki = nakladka.firstElementChild;
  dokument.body.appendChild(nakladka);
  if (korzenNakladki === null) return;
  zwiazNakladke(kanal, korzenNakladki);
  dolozPrzyciskModelu(korzenNakladki);
  await wypelnijTozsamosc(kanal, korzenNakladki);
  await odczytajNastawy(kanal, korzenNakladki);
  await pokazObowiazujace(kanal, korzenNakladki);
  oglos(NAGLOWEK, 'Konfiguracja otwarta.');
}

function zwiazNakladke(kanal: Kanal, nakladka: Element): void {
  nakladka.addEventListener('click', (zdarzenie) => {
    const cel = zdarzenie.target;
    if (!(cel instanceof Element)) return;
    if (cel.closest('[data-zamknij-konfig]') !== null) {
      zdarzenie.stopPropagation();
      nakladka.remove();
      return;
    }
    const nawigacja = cel.closest<HTMLElement>('[data-nav-sekcja]');
    if (nawigacja !== null) {
      zdarzenie.stopPropagation();
      pokazSekcje(nakladka, nawigacja.dataset.navSekcja ?? '');
      return;
    }
    const zapis = cel.closest<HTMLElement>('[data-konfig-zapisz]');
    if (zapis !== null) {
      zdarzenie.stopPropagation();
      void zapiszDokument(kanal, nakladka, zapis.dataset.konfigZapisz ?? '');
      return;
    }
    const zdjecie = cel.closest<HTMLElement>('[data-konfig-zdejmij]');
    if (zdjecie !== null) {
      zdarzenie.stopPropagation();
      void zdejmijDokument(kanal, nakladka, zdjecie.dataset.konfigZdejmij ?? '');
      return;
    }
    if (cel.closest('[data-konfig-przywroc]') !== null) {
      zdarzenie.stopPropagation();
      void przywrocDomyslne(kanal, nakladka);
      return;
    }
    if (cel.closest('[data-konfig-model]') !== null) {
      zdarzenie.stopPropagation();
      void wskazModelWiodacy(kanal, nakladka);
    }
  });
}

function pokazSekcje(nakladka: Element, oznaczenie: string): void {
  for (const sekcja of nakladka.querySelectorAll<HTMLElement>('.kf-sekcja')) {
    sekcja.hidden = sekcja.id !== oznaczenie;
  }
}

function polaSekcji(nakladka: Element, oznaczenie: string): HTMLTextAreaElement | null {
  return nakladka.querySelector<HTMLTextAreaElement>(`#${oznaczenie} .kf-textarea`);
}

async function wypelnijTozsamosc(kanal: Kanal, nakladka: Element): Promise<void> {
  const spis = await wywolaj(kanal, Command.IdentityCategoryList, {});
  if (!spis.udany || spis.wynik === undefined) {
    oglos(NAGLOWEK, spis.blad?.message ?? 'Wykaz kategorii tożsamości nie doszedł.',
      'ostrzezenie');
    return;
  }
  const kategorie = spis.wynik.categories;
  const kategoria = kategorie.find((poz) => poz.id === KOD_KATEGORII);
  const pole = polaSekcji(nakladka, SEKCJA_DOKUMENTU);
  if (kategoria !== undefined && pole !== null) {
    podpiszSekcje(nakladka, SEKCJA_DOKUMENTU, kategoria);
    dolozCzynnosciSekcji(nakladka, SEKCJA_DOKUMENTU, kategoria.id);
    await wczytajDokument(kanal, pole, kategoria.id);
  }
  oglos(NAGLOWEK, `Kategorii tożsamości: ${String(kategorie.length)};`
    + ' w nakładce pisze się Konstytucję, reszta warstw jest do wglądu.');
}

/* Sekcja podpisuje się nazwą kategorii z rdzenia: makieta nosi nazwy własne,
   a obowiązuje ta, pod którą rdzeń prowadzi treść. */
function podpiszSekcje(
  nakladka: Element,
  oznaczenie: string,
  kategoria: IdentityCategory,
): void {
  const etykieta = nakladka.querySelector(`#${oznaczenie} .kf-pole-etyk`);
  if (etykieta !== null) etykieta.textContent = kategoria.name;
  const zasieg = nakladka.querySelector(`#${oznaczenie} .kf-pole-zasieg`);
  if (zasieg !== null) zasieg.textContent = `warstwa: ${kategoria.layer}`;
}

/* Rdzeń oddaje dokumenty kategorii jako zbiór; nakładka pisze pierwszy z nich,
   a jego identyfikator zostaje w polu, bo zdjęcie idzie po dokumencie. */
async function wczytajDokument(
  kanal: Kanal,
  pole: HTMLTextAreaElement,
  idKategorii: string,
): Promise<void> {
  pole.dataset.kategoria = idKategorii;
  const dokument = await wywolaj(kanal, Command.IdentityDocumentGet, {
    categoryId: idKategorii,
  });
  if (!dokument.udany || dokument.wynik === undefined) {
    oglos(NAGLOWEK, dokument.blad?.message ?? 'Treść kategorii nie doszła.',
      'ostrzezenie');
    return;
  }
  const stojacy = dokument.wynik.documents[0];
  pole.value = stojacy?.content ?? '';
  pole.dataset.dokument = stojacy?.id ?? '';
}

function dolozCzynnosciSekcji(
  nakladka: Element,
  oznaczenie: string,
  idKategorii: string,
): void {
  const sekcja = nakladka.querySelector(`#${oznaczenie} .kf-sekcja-cialo`);
  if (sekcja === null || sekcja.querySelector('[data-konfig-zapisz]') !== null) return;
  const dokument = sekcja.ownerDocument;
  const pasek = dokument.createElement('div');
  pasek.className = 'kf-pole-belka';
  pasek.append(
    przyciskCzynnosci(dokument, 'data-konfig-zapisz', idKategorii, 'Zapisz treść',
      'dn-btn dn-btn--sygnal dn-btn--sm'),
    przyciskCzynnosci(dokument, 'data-konfig-zdejmij', idKategorii, 'Zdejmij treść',
      'dn-btn dn-btn--duch dn-btn--sm'),
    przyciskCzynnosci(dokument, 'data-konfig-przywroc', '', 'Przywróć nastawy domyślne',
      'dn-btn dn-btn--duch dn-btn--sm'),
  );
  sekcja.appendChild(pasek);
}

function przyciskCzynnosci(
  dokument: Document,
  klucz: string,
  wartosc: string,
  napis: string,
  klasa: string,
): HTMLButtonElement {
  const przycisk = dokument.createElement('button');
  przycisk.type = 'button';
  przycisk.className = klasa;
  przycisk.setAttribute(klucz, wartosc);
  przycisk.textContent = napis;
  return przycisk;
}

async function zapiszDokument(
  kanal: Kanal,
  nakladka: Element,
  idKategorii: string,
): Promise<void> {
  const pole = nakladka.querySelector<HTMLTextAreaElement>(
    `.kf-textarea[data-kategoria="${idKategorii}"]`);
  if (pole === null) {
    oglos(NAGLOWEK, 'Nakładka nie ma pola tej kategorii.', 'ostrzezenie');
    return;
  }
  const wynik = await wywolaj(kanal, Command.IdentityDocumentSet, {
    categoryId: idKategorii,
    content: pole.value,
    enabled: true,
  });
  if (!wynik.udany) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił zapisu dokumentu.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, 'Dokument tożsamości zapisany.');
  await pokazObowiazujace(kanal, nakladka);
}

async function zdejmijDokument(
  kanal: Kanal,
  nakladka: Element,
  idKategorii: string,
): Promise<void> {
  const pole = nakladka.querySelector<HTMLTextAreaElement>(
    `.kf-textarea[data-kategoria="${idKategorii}"]`);
  const idDokumentu = pole?.dataset.dokument ?? '';
  if (idDokumentu === '') {
    oglos(NAGLOWEK, 'Ta kategoria nie ma jeszcze zapisanej treści.', 'ostrzezenie');
    return;
  }
  const wynik = await wywolaj(kanal, Command.IdentityDocumentRemove, {
    documentId: idDokumentu,
  });
  if (!wynik.udany) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił zdjęcia dokumentu.', 'ostrzezenie');
    return;
  }
  if (pole !== null) {
    pole.value = '';
    pole.dataset.dokument = '';
  }
  oglos(NAGLOWEK, 'Dokument tożsamości zdjęty; warstwa przestaje obowiązywać.');
  await pokazObowiazujace(kanal, nakladka);
}

/* Podgląd mówi, co obowiązuje po złożeniu wszystkich warstw — to jedyne miejsce,
   w którym widać wynik, a nie samą treść wpisaną w polu. */
async function pokazObowiazujace(kanal: Kanal, nakladka: Element): Promise<void> {
  const podglad = nakladka.querySelector('.kf-podglad');
  if (podglad === null) return;
  const tozsamosc = await wywolaj(kanal, Command.IdentityEffectiveGet, {});
  const nastawy = await wywolaj(kanal, Command.ConfigEffectiveGet, {
    scope: ConfigScope.Application,
  });
  if (tozsamosc.udany && tozsamosc.wynik !== undefined) {
    pokazWarstwy(nakladka, tozsamosc.wynik.layers);
  }
  const dokument = podglad.ownerDocument;
  podglad.replaceChildren();
  const wiersze = [
    `tożsamość obowiązująca: ${tozsamosc.udany
      ? `${String(tozsamosc.wynik?.layers.length ?? 0)} warstw`
        + ` · ${String(tozsamosc.wynik?.prompt.length ?? 0)} znaków`
      : (tozsamosc.blad?.message ?? 'bez odpowiedzi')}`,
    `nastawy rozstrzygnięte: ${nastawy.udany
      ? String(nastawy.wynik?.effective.origins.length ?? 0) + ' źródeł'
      : (nastawy.blad?.message ?? 'bez odpowiedzi')}`,
  ];
  for (const tresc of wiersze) {
    const wiersz = dokument.createElement('div');
    wiersz.textContent = tresc;
    podglad.appendChild(wiersz);
  }
}

/* Przywrócenie zdejmuje nastawy z zasięgu aplikacji: sesje i okna zachowują
   własne, węższe nastawy, bo te stoją w innym zasięgu. */
async function przywrocDomyslne(kanal: Kanal, nakladka: Element): Promise<void> {
  const wynik = await wywolaj(kanal, Command.ConfigReset, {
    scope: ConfigScope.Application,
  });
  if (!wynik.udany) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił przywrócenia nastaw.',
      'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, 'Nastawy zasięgu aplikacji przywrócone do domyślnych.');
  await odczytajNastawy(kanal, nakladka);
  await pokazObowiazujace(kanal, nakladka);
}

/* Nastawy sesji czyta się osobno od tożsamości: to dwie różne warstwy, które
   spotykają się dopiero w podglądzie obowiązującym. */
async function odczytajNastawy(kanal: Kanal, nakladka: Element): Promise<void> {
  const wynik = await wywolaj(kanal, Command.ConfigSessionGet, {
    scope: ConfigScope.Application,
    areas: [SessionConfigArea.Model, SessionConfigArea.Permissions],
  });
  if (!wynik.udany) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Nastawy sesji nie doszły.', 'ostrzezenie');
    return;
  }
  const wybor = nakladka.querySelector<HTMLSelectElement>('select[data-model]');
  if (wybor === null) return;
  const model = wynik.wynik?.config.model?.primaryModel ?? '';
  /* Nastawa z rdzenia bywa spoza wykazu makiety, więc wchodzi jako własna
     pozycja — inaczej pole pokazywałoby model, którego nikt nie wskazał. */
  if (model !== '' && !Array.from(wybor.options).some((poz) => poz.value === model)) {
    const pozycja = wybor.ownerDocument.createElement('option');
    pozycja.value = model;
    pozycja.textContent = `${model} (z nastaw)`;
    wybor.appendChild(pozycja);
  }
  wybor.value = model;
}

/* Profil i Ekspertyza są w makiecie tylko do odczytu, więc dostają to, co
   rdzeń podaje jako warstwy obowiązujące — inaczej stałyby na treści makiety. */
function pokazWarstwy(nakladka: Element, warstwy: readonly IdentityLayerContent[]): void {
  const gniazda: ReadonlyArray<readonly [string, number]> = [
    ['kf-sek-profil', 1],
    ['kf-sek-ekspertyza', 2],
  ];
  for (const [oznaczenie, numer] of gniazda) {
    const podglad = nakladka.querySelector(`#${oznaczenie} .kf-podglad`);
    if (podglad === null) continue;
    const warstwa = warstwy[numer];
    podglad.textContent = warstwa === undefined
      ? 'Rdzeń nie prowadzi treści dla tej warstwy.'
      : `${warstwa.name}: ${String(warstwa.content.length)} znaków`
        + ` · oś ${warstwa.axis}`;
  }
}

/* Makieta pokazuje wykaz modeli, ale niczego nim nie zapisuje: przycisk obok
   wykazu domyka drogę od wskazania do nastawy w rdzeniu. */
function dolozPrzyciskModelu(nakladka: Element): void {
  const wybor = nakladka.querySelector('select[data-model]');
  if (wybor === null) return;
  const przycisk = nakladka.ownerDocument.createElement('button');
  przycisk.type = 'button';
  przycisk.className = 'dn-btn dn-btn--atrament dn-btn--sm';
  przycisk.dataset.konfigModel = 'wiodacy';
  przycisk.textContent = 'Zapisz model wiodący';
  wybor.after(przycisk);
}

/* Jednostką zapisu jest obszar, nie pole: zapis obszaru modelu zostawia
   uprawnienia nietknięte, więc pytanie obejmuje wyłącznie model wiodący. */
async function wskazModelWiodacy(kanal: Kanal, nakladka: Element): Promise<void> {
  const wybor = nakladka.querySelector<HTMLSelectElement>('select[data-model]');
  const odpowiedz = await zapytajWOknieModalnym({
    tytul: 'Model wiodący zasięgu aplikacji',
    opis: 'Wskazanie obowiązuje sesje, które nie mają własnego modelu.',
    pola: [{
      klucz: 'model',
      etykieta: 'Oznaczenie modelu',
      wartosc: wybor?.value ?? '',
      wymagane: true,
    }],
    wykonanie: 'Zapisz nastawę',
  });
  if (odpowiedz === null) {
    oglos(NAGLOWEK, 'Wskazanie modelu odwołane; nastawy bez zmiany.');
    return;
  }
  const model = (odpowiedz['model'] ?? '').trim();
  if (model === '') {
    oglos(NAGLOWEK, 'Puste oznaczenie nie jest wskazaniem.', 'ostrzezenie');
    return;
  }
  const wynik = await wywolaj(kanal, Command.ConfigSessionSet, {
    scope: ConfigScope.Application,
    areas: [SessionConfigArea.Model],
    config: { model: { primaryModel: model } },
  });
  if (!wynik.udany) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił zapisu nastawy.',
      'ostrzezenie');
    return;
  }
  const pominiete = wynik.wynik?.unsupportedFields ?? [];
  oglos(NAGLOWEK, 'Nastawa modelu zapisana; obszarów: '
    + String(wynik.wynik?.storedAreas.length ?? 0)
    + (pominiete.length > 0
      ? `, pól bez obsługi dostawcy: ${String(pominiete.length)}.`
      : '.'));
  await odczytajNastawy(kanal, nakladka);
  await pokazObowiazujace(kanal, nakladka);
}

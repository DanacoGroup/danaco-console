import { AppArchitectureTemplate as SzablonArchitektury } from '../../../../shared/contract';
import type {
  AppArchitecture,
  AppArchitectureTemplate,
  AppComponent,
} from '../../../../shared/contract';
import { opisOdmowy } from '../../komponenty/odmowa';
import { poleTekstowe, przycisk, utworzWierszOdpowiedzi } from '../../modele/kontrolki-formularza';
import { utworzWykazBrakow } from './braki-kontraktu';
import { opiszPole } from './dymek-objasnienia';
import { BRAKI_ARCHITECTURE_DESIGNER, KODY_OKIEN, NAZWY_OKIEN, BEZ_OKNA_MODULU } from './etykiety-apps';
import { utworzFormularzKomponentu } from './formularz-komponentu';
import { utworzKanweKomponentow } from './kanwa-komponentow';
import { utworzRameApps } from './rama-okna';
import { narzedziaArchitectureDesigner } from './narzedzia-apps';
import { utworzPrzybornikApps } from './przybornik-apps';
import type { StanProduktu } from './stan-produktu';
import { utworzWyborZMenu, wierszWyboru } from './wybor-z-menu';

/**
 * Architecture Designer — okno kreatora modułu Apps, punkt wyjścia procesu.
 *
 * Definiowanie komponentów rozwiązania i ustalenie zależności stoją na jednej
 * komendzie kontraktu — `apps.architecture.define` — która przyjmuje szablon
 * i komplet komponentów wraz z ich `dependsOn`.
 *
 * „Waliduj” to ta sama komenda, nie druga droga: kontrakt nie ma osobnej
 * komendy walidacji, zwraca za to `validationIssues` w architekturze po zapisie.
 * Przycisk pozostaje klikalny przy pustej kanwie — odpowiedź jest opisowa,
 * nie blokadą.
 *
 * Zdanie potwierdzające bierze treść z odpowiedzi, nie z zamówienia: wymienia
 * identyfikator, wersję, szablon i komponenty tak, jak oddał je rdzeń, a każdą
 * rozbieżność z zamówieniem i każde zastrzeżenie walidacji ogłasza odmową.
 */
export interface OknoArchitectureDesigner {
  element: HTMLElement;
  odswiez(): void;
}

export function utworzOknoArchitectureDesigner(stan: StanProduktu): OknoArchitectureDesigner {
  const kod = KODY_OKIEN.ArchitectureDesigner;
  const rama = utworzRameApps(kod, NAZWY_OKIEN[kod] ?? kod, 'kreator');

  const nazwa = poleTekstowe({ etykieta: 'Nazwa architektury', podpowiedz: 'np. Wydanie 1' });
  // Rozwijanie z biblioteki kontrolek, nie natywny `<select>`.
  const szablon = utworzWyborZMenu('Szablon architektury', [
    { wartosc: SzablonArchitektury.Monolith, etykieta: 'Monolit' },
    { wartosc: SzablonArchitektury.Microservices, etykieta: 'Mikroserwisy' },
    { wartosc: SzablonArchitektury.Serverless, etykieta: 'Architektura bezserwerowa' },
  ]);
  const szablonWiersz = wierszWyboru('Szablon architektury', szablon);

  const formularz = utworzFormularzKomponentu();
  const kanwa = utworzKanweKomponentow((identyfikator) => stan.usunKomponent(identyfikator));

  const dodaj = przycisk('Dodaj komponent do kanwy', 'dn-btn dn-btn--sm dn-btn--zarys');
  const waliduj = przycisk('Waliduj', 'dn-btn dn-btn--sm dn-btn--zarys');
  const zapisz = przycisk('Zapisz architekturę', 'dn-btn dn-btn--sm dn-btn--atrament');
  // Odczyt stoi w tym samym pasku co zapis, bo jest jego odwrotnością. Bez
  // niego kanwa zaczyna pusta przy każdym wejściu w moduł, choć rdzeń trzyma
  // architekturę w bazie.
  const odczytaj = przycisk('Odczytaj architekturę z rdzenia', 'dn-btn dn-btn--sm dn-btn--zarys');
  const odpowiedz = utworzWierszOdpowiedzi();

  const wersje = document.createElement('p');
  wersje.className = 'mp-wersje';

  const pasek = document.createElement('div');
  pasek.className = 'mp-pasek';
  pasek.append(dodaj, waliduj, zapisz, odczytaj);

  rama.akcje.append(utworzWykazBrakow('Bez drogi w kontrakcie', BRAKI_ARCHITECTURE_DESIGNER));
  rama.akcje.append(
    utworzPrzybornikApps('Walidacja, wersje, adnotacje i eksport', narzedziaArchitectureDesigner(stan))
      .element,
  );
  rama.tresc.append(
    opiszPole(nazwa.element, 'Nazwa architektury bywa pusta — kontrakt jej nie wymaga.'),
    opiszPole(szablonWiersz, 'Szablon z wyliczenia kontraktu: monolit, mikroserwisy, bezserwerowa.'),
    formularz.element,
    pasek,
    kanwa.element,
    wersje,
    odpowiedz.element,
  );

  dodaj.addEventListener('click', () => {
    const komponent = formularz.zbierz();
    if (komponent === null) {
      odpowiedz.pokaz('Wskaż nazwę komponentu — bez niej kanwa nie ma czego przyjąć.', false);
      return;
    }
    stan.dodajKomponent(komponent);
    formularz.wyczysc();
    odpowiedz.pokaz(`Komponent „${komponent.name}" zestawiony na kanwie.`, true);
  });

  waliduj.addEventListener('click', () => void wyslij('Walidacja układu'));
  zapisz.addEventListener('click', () => void wyslij('Zapis architektury'));
  odczytaj.addEventListener('click', () => void odczytajArchitekture());

  /**
   * Odczyt architektury zapisanej w rdzeniu.
   *
   * Odczyt nadpisuje kanwę i mówi o tym wprost: architektura oddana przez rdzeń
   * wchodzi wraz ze swoim kompletem komponentów, więc komponenty zestawione
   * na kanwie, a jeszcze niezapisane, po tym kroku znikają. Zdanie wymienia
   * liczbę komponentów, która przyszła.
   *
   * Brak architektury jest odpowiedzią, nie odmową — rdzeń oddaje wtedy wynik
   * udany bez pola, kanwa zostaje nietknięta.
   */
  async function odczytajArchitekture(): Promise<void> {
    rama.ladowanie('Odczyt architektury z rdzenia…');
    odpowiedz.pokaz('Odczyt architektury: żądanie wysłane do rdzenia…', true);
    const byla = await stan.odczytajArchitekture();
    const powod = stan.powodOdczytu('architektura');
    if (powod !== '') {
      rama.blad(powod);
      odpowiedz.pokaz(powod, false);
      return;
    }
    const oddana = stan.architektura();
    rama.gotowe();
    odswiez();
    if (!byla || oddana === null) {
      odpowiedz.pokaz(
        'Odczyt architektury: rdzeń nie zna żadnej architektury tego okna — kanwa została ' +
          'bez zmiany. Zdefiniuj komponenty i naciśnij „Zapisz architekturę".',
        true,
      );
      return;
    }
    nazwa.kontrolka.value = oddana.name ?? '';
    // Szablon spoza wyliczenia nie wchodzi na uchwyt po cichu: wykaz niesie
    // trzy wartości kontraktu, a szablon czwarty nie ma w nim pozycji do
    // wskazania. Zdanie mówi wtedy o rozejściu okna z rdzeniem.
    const szablonPrzyjety = szablon.ustawWartosc(oddana.template);
    const uwagaSzablonu = szablonPrzyjety
      ? ''
      : ` UWAGA: rdzeń oddał szablon ${oddana.template}, którego nie ma w wyliczeniu kontraktu — ` +
        'wybór w oknie został przy poprzedniej wartości.';
    odpowiedz.pokaz(
      `Odczyt architektury: rdzeń oddał układ ${oddana.id} w wersji ` +
        `${oddana.version ?? 'bez numeru'} — szablon ${oddana.template}, komponentów ` +
        `${stan.komponenty().length}. Kanwa pokazuje teraz stan z rdzenia.${uwagaSzablonu}`,
      szablonPrzyjety,
    );
  }

  async function wyslij(czynnosc: string): Promise<void> {
    const idOkna = stan.idOkna();
    if (idOkna === '') {
      odpowiedz.pokaz(BEZ_OKNA_MODULU, false);
      return;
    }
    const zamowione = stan.komponenty();
    const zamowionySzablon = szablon.wartosc() as AppArchitectureTemplate;
    rama.ladowanie(`${czynnosc} w toku…`);
    odpowiedz.pokaz(`${czynnosc}: żądanie wysłane do rdzenia…`, true);
    const wynik = await stan.zrodlo.zapiszArchitekture({
      idOkna,
      idArchitektury: stan.architektura()?.id ?? '',
      nazwa: nazwa.kontrolka.value,
      szablon: zamowionySzablon,
      komponenty: zamowione,
    });
    if (!wynik.udany || wynik.wynik === undefined) {
      rama.blad(opisOdmowy(czynnosc, wynik.blad?.code, wynik.blad?.message));
      odpowiedz.pokaz(opisOdmowy(czynnosc, wynik.blad?.code, wynik.blad?.message), false);
      return;
    }
    const oddana = wynik.wynik.architecture;
    stan.wchlonArchitekture(oddana);
    const rozbieznosc = rozbieznoscUkladu(zamowione, zamowionySzablon, oddana);
    if (rozbieznosc !== '') {
      const zdanie =
        `${czynnosc}: rdzeń zapisał układ ${oddana.id}, ale ODDAŁ CO INNEGO, NIŻ ` +
        `ZAMÓWIONO — ${rozbieznosc}. Zapisu nie potwierdzam; na kanwie stoi teraz ` +
        'układ oddany przez rdzeń.';
      rama.blad(zdanie);
      odpowiedz.pokaz(zdanie, false);
      return;
    }
    // Przesłona zdejmowana tutaj, nie w `odswiez`: `odswiez` wychodzi przy
    // fazie `ladowanie`, więc bez zamknięcia fazy okno pokazywałoby wskaźnik
    // pracy po jej zakończeniu. Po zamknięciu fazy przeliczamy widok, bo
    // wchłonięcie architektury ogłosiło zmianę jeszcze w fazie `ladowanie`
    // i tamten przebieg wyszedł bez rozstrzygnięcia między pustą kanwą
    // a kanwą z komponentami.
    rama.gotowe();
    odswiez();
    const zastrzezenia = oddana.validationIssues ?? [];
    const opisUkladu =
      `układ ${oddana.id} w wersji ${oddana.version ?? 'bez numeru'} — szablon ` +
      `${oddana.template}, komponentów ${(oddana.components ?? []).length}`;
    if (zastrzezenia.length > 0) {
      // Zastrzeżenie rdzenia nie jest potwierdzeniem — zdanie musi je wymienić,
      // tak samo jak robi to wiersz „Wersje” pod spodem.
      odpowiedz.pokaz(
        `${czynnosc}: rdzeń zapisał ${opisUkladu}, ale ZGŁOSIŁ ZASTRZEŻENIA ` +
          `(${zastrzezenia.length}): ${zastrzezenia.join(' · ')}`,
        false,
      );
      return;
    }
    odpowiedz.pokaz(
      `${czynnosc}: rdzeń oddał ${opisUkladu}. ${opisPustychZastrzezen(oddana)}`,
      true,
    );
  }

  function odswiez(): void {
    const komponenty = stan.komponenty();
    kanwa.nanies(komponenty);
    formularz.ustawZaleznosci(komponenty);
    wersje.textContent = opisWersji(stan);
    if (rama.faza() === 'blad' || rama.faza() === 'ladowanie') return;
    if (komponenty.length === 0) {
      rama.puste('Kanwa bez zestawionych jeszcze komponentów.');
      return;
    }
    rama.gotowe();
  }

  odswiez();
  return { element: rama.element, odswiez };
}

/**
 * Czym oddany układ różni się od zamówionego — pusty łańcuch, gdy niczym.
 *
 * Zestawiane są zbiory identyfikatorów, nie listy. Rdzeń trzyma komponenty
 * i krawędzie zależności w osobnych tabelach i oddaje je w porządku własnym,
 * a powtórzoną krawędź zdejmuje kluczem pierwotnym. Porządek i powtórzenie nie
 * są więc rozbieżnością, tylko cudzym sposobem przechowywania; rozbieżnością
 * jest komponent, który zniknął, albo taki, którego nikt nie zamawiał.
 */
function rozbieznoscUkladu(
  zamowione: readonly AppComponent[],
  zamowionySzablon: string,
  oddana: AppArchitecture,
): string {
  const rozejscia: string[] = [];
  if (oddana.template !== zamowionySzablon) {
    rozejscia.push(`zamówiono szablon ${zamowionySzablon}, rdzeń oddał ${oddana.template}`);
  }
  const oddane = new Set((oddana.components ?? []).map((komponent) => komponent.id));
  const brakujace = zamowione
    .filter((komponent) => !oddane.has(komponent.id))
    .map((komponent) => komponent.id);
  if (brakujace.length > 0) {
    rozejscia.push(`rdzeń nie oddał komponentów: ${brakujace.join(', ')}`);
  }
  const zamowioneKody = new Set(zamowione.map((komponent) => komponent.id));
  const nadmiarowe = [...oddane].filter((identyfikator) => !zamowioneKody.has(identyfikator));
  if (nadmiarowe.length > 0) {
    rozejscia.push(`rdzeń oddał komponenty spoza zamówienia: ${nadmiarowe.join(', ')}`);
  }
  const zgubioneKrawedzie = zgubioneZaleznosci(zamowione, oddana.components ?? []);
  if (zgubioneKrawedzie.length > 0) {
    rozejscia.push(`rdzeń nie oddał zależności: ${zgubioneKrawedzie.join(', ')}`);
  }
  return rozejscia.join('; ');
}

/**
 * Które zamówione krawędzie zależności nie wróciły z rdzenia.
 *
 * Zależność jest wskazaniem identyfikatora innego komponentu, więc jej zgubienie
 * zmienia układ tak samo jak zgubienie komponentu, a kanwa zaciągnięta
 * z odpowiedzi sama z siebie nie pokazuje, że zamówiono więcej.
 */
function zgubioneZaleznosci(
  zamowione: readonly AppComponent[],
  oddane: readonly AppComponent[],
): string[] {
  const wedlugKodu = new Map(oddane.map((komponent) => [komponent.id, komponent]));
  const zgubione: string[] = [];
  for (const komponent of zamowione) {
    const oddany = wedlugKodu.get(komponent.id);
    if (oddany === undefined) continue;
    const maja = new Set(oddany.dependsOn ?? []);
    for (const zrodlo of komponent.dependsOn ?? []) {
      if (!maja.has(zrodlo)) zgubione.push(`${zrodlo} → ${komponent.id}`);
    }
  }
  return zgubione;
}

/**
 * Zdanie o pustym polu `validationIssues` — orzeka o odebranej ramce, nie
 * o zachowaniu rdzenia.
 *
 * Rdzeń tego pola nie wylicza: adapter czyta zastrzeżenia z bazy, zapisuje je
 * z powrotem tym samym UPSERT-em i oddaje w odpowiedzi, więc jedynym pisarzem
 * kolumny jest ten, kto ją przed chwilą przeczytał. Zdanie „bez zastrzeżeń"
 * byłoby zapewnieniem o sprawdzeniu, którego nikt nie wykonał.
 *
 * Własny walidator w oknie byłby drugą prawdą o tym samym układzie i zacierał
 * granicę odpowiedzialności: zastrzeżenia widoczne w oknie sugerowałyby, że zna
 * je rdzeń.
 *
 * Rozróżnienie pola nieobecnego od pola obecnego i pustego sprawia, że w dniu,
 * w którym rdzeń zacznie zastrzeżenia liczyć, zdanie zmieni treść bez tknięcia
 * tego pliku.
 */
function opisPustychZastrzezen(architektura: AppArchitecture): string {
  // Kontrakt deklaruje pole jako `?: string[]`, ale gniazdo przynosi w tym
  // miejscu JSON-owy `null`. Rzutowanie nie poszerza umowy — nazywa kształt,
  // który naprawdę przychodzi, żeby rozróżnienie „pola nie ma" od „pole jest
  // puste" opierało się na ramce, a nie na deklaracji.
  const pole = architektura.validationIssues as readonly string[] | null | undefined;
  if (pole === undefined || pole === null) {
    return (
      'Pole validationIssues nie przyszło w tej odpowiedzi — to nie jest ' +
      'orzeczenie, że układ jest poprawny.'
    );
  }
  return 'Pole validationIssues przyszło puste.';
}

/** Wiersz „Wersje” panelu akcji — numer i zastrzeżenia z odpowiedzi rdzenia. */
function opisWersji(stan: StanProduktu): string {
  const architektura = stan.architektura();
  if (architektura === null) {
    return 'Wersje: rdzeń nie potwierdził jeszcze żadnego zapisu architektury.';
  }
  const zastrzezenia = architektura.validationIssues ?? [];
  const numer = architektura.version ?? 0;
  if (zastrzezenia.length === 0) {
    return `Wersja ${numer} — ${opisPustychZastrzezen(architektura)}`;
  }
  return `Wersja ${numer} — zastrzeżenia rdzenia: ${zastrzezenia.join(' · ')}`;
}

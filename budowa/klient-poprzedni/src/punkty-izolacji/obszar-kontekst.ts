import {
  Command,
  IsolationContextKind,
  IsolationLayer,
  type IsolationSwitch,
} from '../../../shared/contract';
import { wiersz, wybor } from '../modele/kontrolki-formularza-braki';
import { utworzDymekObjasnienia } from '../komponenty/dymek';
import { opisOdmowyBledu } from '../komponenty/odmowa';
import { utworzStanTresci } from '../komponenty/stan-tresci';
import type { Kanal, Wynik } from '../protokol/kanal';
import type { ObszarIzolacji, ZaleznosciObszaru } from './obszary';
import { NAZWY_WARSTW } from './stan-warstwy';

/**
 * Obszar „Kontekst” — trzy punkty izolacji kontekstu rozstrzygane niezależnie
 * od siebie (`server/internal/konfig/definicje_izolacji.go`,
 * `definicjeIzolacjiKontekstu`, `IsolationContextKind.History/Memory/Context`):
 * historia wymiany, pamięć długoterminowa i bieżący stan roboczy. Każdy
 * przyjmuje jedną z dwóch wartości — `odrebna` albo `wspoldzielona` — z
 * domyślną `odrebna` (rdzeń startuje z pełną izolacją; współdzielenie jest
 * zawsze decyzją Operatora, nigdy stanem narzuconym).
 *
 * Treść wyjaśnień „odrębna”/„współdzielona” przy każdym kluczu jest przepisana
 * z `Definicja.Objasnienie` tego samego pliku rdzenia — jedno źródło prawdy,
 * żeby zdanie widziane przez Operatora nie rozjechało się z tym, co egzekwuje
 * maszyneria.
 *
 * Obszar woła `isolation.context.get` i `isolation.context.set` wprost;
 * `kanal.wyslij` jest opakowany w obietnicę, a odmowa rozstrzygana przez
 * `opisOdmowyBledu`.
 *
 * Poziom zasięgu i warstwa nie są tu zaszyte: oba przychodzą ze stanu wspólnego
 * okna (`stan-zasiegu.ts`, `stan-warstwy.ts`). Zaszycie ich znaczyłoby, że
 * selektor zasięgu w lewym panelu i wybór warstwy w pasie narzędzi pokazują
 * jedno, a odczyt idzie po drugie.
 *
 * „Wartość obecna” ma trzy stany, celowo nierysowane identycznie:
 *   - zapisana — odczyt na warstwie czynnej powiódł się: Operator świadomie
 *     rozstrzygnął tę wartość na tym zasięgu i tej warstwie;
 *   - domyślna — odczyt warstwy czynnej zwrócił `not_found`: zapisu Operatora
 *     nie ma, więc widok pokazuje wartość domyślną platformy jawnie oznaczoną
 *     jako domyślna, a nie jako świadomy zapis;
 *   - nieodczytana — odczyt zwrócił inną odmowę niż `not_found` (usterka
 *     rdzenia, brak uprawnienia, kanał niedostępny): to nie jest informacja
 *     o stanie maszyny, tylko brak informacji, i tak też jest nazwana.
 *
 * Żaden z tych stanów nie jest rysowany jako sukces i nie ma tu zapisu
 * optymistycznego: widok zmienia się dopiero po odpowiedzi
 * `isolation.context.set`, nigdy przed nią. Selekt wyboru nowej wartości
 * i przycisk zapisu zostają zawsze czynne; odmowa rdzenia jest meldowana
 * zdaniem, nie blokadą kontrolki.
 */

const WARTOSC_ODREBNA = 'odrebna';
const WARTOSC_WSPOLDZIELONA = 'wspoldzielona';

const OPCJE_WARTOSCI: ReadonlyArray<readonly [string, string]> = [
  [WARTOSC_ODREBNA, 'Odrębna'],
  [WARTOSC_WSPOLDZIELONA, 'Współdzielona'],
];

interface PunktKontekstu {
  /** Rodzaj kontekstu w kontrakcie — `IsolationContextKind`. */
  kind: IsolationContextKind;
  /** Klucz konfiguracji rdzenia, wyłącznie do wyświetlenia obok nazwy. */
  klucz: string;
  /** Nazwa czytelna po polsku, widoczna w nagłówku karty. */
  nazwa: string;
  /** Wyjaśnienie skutku obu wartości — Operator ma rozumieć skutek, nie samą nazwę. */
  objasnienie: string;
}

const PUNKTY: readonly PunktKontekstu[] = [
  {
    kind: IsolationContextKind.History,
    klucz: 'izolacja_historia',
    nazwa: 'Historia wymiany',
    objasnienie:
      'Zapis wymiany wiadomości. Odrębna: nowe okno zaczyna z pustą historią. ' +
      'Współdzielona: ten sam zapis widoczny w kilku oknach lub modułach.',
  },
  {
    kind: IsolationContextKind.Memory,
    klucz: 'izolacja_pamiec',
    nazwa: 'Pamięć długoterminowa',
    objasnienie:
      'Pamięć długoterminowa zasięgu. Odrębna: pamięć jednego zasięgu niewidoczna ' +
      'w innym. Współdzielona: jeden zasób zasila kilka zasięgów.',
  },
  {
    kind: IsolationContextKind.Context,
    klucz: 'izolacja_kontekst',
    nazwa: 'Stan roboczy (kontekst bieżący)',
    objasnienie:
      'Bieżący stan roboczy: aktywne pliki, projekt, załączniki, zmienne. Odrębny: ' +
      'nie przenosi się między oknami, każde zaczyna od nowa. Współdzielony przenosi ' +
      'się między oknami bez przeładowania.',
  },
];

/** Rozstrzygnięcie „wartości obecnej” jednego punktu po odczycie warstwy sesji. */
type StanObecny =
  | { rodzaj: 'zapisane'; isolated: boolean }
  | { rodzaj: 'domyslne' }
  | { rodzaj: 'blad'; opis: string };

/** Opakowuje `kanal.wyslij` w Promise — kanał sam daje wyłącznie wersję z wywołaniem zwrotnym. */
function posijKomende<T>(
  kanal: Kanal,
  komenda: Command,
  zadanie: unknown,
): Promise<Wynik<T>> {
  return new Promise((rozstrzygnij) => {
    kanal.wyslij(komenda as never, zadanie as never, (wynik) => rozstrzygnij(wynik as Wynik<T>));
  });
}

export function utworzObszar(zaleznosci: ZaleznosciObszaru): ObszarIzolacji {
  const { kanal, warstwa, zasieg } = zaleznosci;
  const tresc = utworzStanTresci('pi');

  /** Ostatnio znane wartości trzech kluczy (do złożenia żądania zapisu pełnym zestawem). */
  const znane = new Map<string, boolean>();

  function odswiez(): void {
    void wczytaj();
  }

  async function wczytaj(): Promise<void> {
    tresc.ladowanie('Pytam rdzeń o punkty izolacji kontekstu (isolation.context.get)…');

    const domyslne = await posijKomende<{ switches: IsolationSwitch[] }>(
      kanal,
      Command.IsolationContextGet,
      { scope: zasieg.zasieg(), scopeId: zasieg.bytDoZadania(), layer: IsolationLayer.Default },
    );
    if (!domyslne.udany || domyslne.wynik === undefined) {
      tresc.blad(
        opisOdmowyBledu('Odczyt wartości domyślnych kontekstu', domyslne.blad),
        domyslne.blad,
      );
      return;
    }

    const obecne = await posijKomende<{ switches: IsolationSwitch[] }>(
      kanal,
      Command.IsolationContextGet,
      { scope: zasieg.zasieg(), scopeId: zasieg.bytDoZadania(), layer: warstwa.warstwa() },
    );

    const stanyObecne = new Map<string, StanObecny>();
    if (obecne.udany && obecne.wynik !== undefined) {
      for (const p of obecne.wynik.switches) stanyObecne.set(p.kind, { rodzaj: 'zapisane', isolated: p.isolated });
    } else if (obecne.blad?.code === 'not_found') {
      // Nikt nie zapisał na warstwie karty sesji — to nie jest usterka odczytu,
      // to brak decyzji Operatora. Pozostaje bez wpisu w `stanyObecne`, każdy
      // punkt dostaje domyślnie rozdzielczość „domyslne” niżej.
    } else {
      const opis = opisOdmowyBledu('Odczyt wartości obecnej', obecne.blad);
      for (const punkt of PUNKTY) stanyObecne.set(punkt.kind, { rodzaj: 'blad', opis });
    }

    znane.clear();
    for (const p of domyslne.wynik.switches) znane.set(p.kind, p.isolated);
    for (const [kind, stan] of stanyObecne) {
      if (stan.rodzaj === 'zapisane') znane.set(kind, stan.isolated);
    }

    const wstep = document.createElement('p');
    wstep.className = 'dn-pole-opis';
    wstep.textContent =
      'Trzy punkty izolacji kontekstu, rozstrzygane niezależnie. Odczyt i zapis idą na zasięgu ' +
      `wskazanym w lewym panelu (${zasieg.opis()}) i na warstwie „${NAZWY_WARSTW[warstwa.warstwa()]}”. ` +
      'Stan wyjściowy platformy: wszystkie trzy odrębne — pełna izolacja, zero współdzielenia bez ' +
      'decyzji Operatora.';

    const kontener = document.createElement('div');
    kontener.className = 'pi-kontekst-karty pi-lista';
    for (const punkt of PUNKTY) {
      const domyslnyPunkt = domyslne.wynik.switches.find((p) => p.kind === punkt.kind);
      const stanObecny = stanyObecne.get(punkt.kind) ?? { rodzaj: 'domyslne' as const };
      kontener.append(zbudujKarte(punkt, domyslnyPunkt?.isolated ?? true, stanObecny));
    }

    const miejsce = tresc.tresc();
    miejsce.append(wstep, kontener);
  }

  function zbudujKarte(
    punkt: PunktKontekstu,
    domyslnyIsolated: boolean,
    stanObecny: StanObecny,
  ): HTMLElement {
    const karta = document.createElement('article');
    karta.className = 'dn-karta';

    const naglowek = document.createElement('header');
    naglowek.className = 'dn-karta-naglowek';
    const tytul = document.createElement('h3');
    tytul.className = 'dn-karta-tytul';
    tytul.textContent = punkt.nazwa;
    const kod = document.createElement('code');
    kod.textContent = punkt.klucz;
    naglowek.append(tytul, kod, utworzDymekObjasnienia(objasnienieKontekstowe(punkt)));

    const cialo = document.createElement('div');
    cialo.className = 'dn-karta-cialo pi-karta__cialo';

    const objasnienie = document.createElement('p');
    objasnienie.className = 'dn-pole-opis';
    objasnienie.textContent = punkt.objasnienie;

    cialo.append(
      objasnienie,
      wierszWartosci(
        'Wartość domyślna platformy',
        plakietka(etykietaWartosci(domyslnyIsolated), 'informacja'),
      ),
      wierszWartosciObecnej(stanObecny),
      zbudujWierszZmiany(punkt, wartoscPoczatkowaSelektu(stanObecny, domyslnyIsolated)),
    );

    karta.append(naglowek, cialo);
    return karta;
  }

  function wierszWartosciObecnej(stan: StanObecny): HTMLElement {
    if (stan.rodzaj === 'zapisane') {
      return wierszWartosci(
        'Wartość obecna',
        plakietka(etykietaWartosci(stan.isolated), 'informacja'),
        `Zapisana świadomie przez Operatora — ${zasieg.opis()}, warstwa „${NAZWY_WARSTW[warstwa.warstwa()]}”.`,
      );
    }
    if (stan.rodzaj === 'domyslne') {
      return wierszWartosci(
        'Wartość obecna',
        plakietka('Domyślna platformy — brak zapisu Operatora', 'sygnal'),
        `Nikt nie zapisał tego punktu na tym zasięgu (${zasieg.opis()}) i warstwie ` +
          `„${NAZWY_WARSTW[warstwa.warstwa()]}”; pokazana wartość pochodzi z domyślnej konfiguracji ` +
          'platformy, nie jest świadomym rozstrzygnięciem.',
      );
    }
    return wierszWartosci(
      'Wartość obecna',
      plakietka('Nie udało się odczytać', 'blad'),
      stan.opis,
    );
  }

  /** Wiersz wyłącznie do odczytu: etykieta, wartość (zwykle plakietka), opcjonalny opis. */
  function wierszWartosci(etykieta: string, wartosc: HTMLElement, opis?: string): HTMLElement {
    const element = document.createElement('div');
    element.className = 'dn-pole';

    const podpis = document.createElement('span');
    podpis.className = 'dn-pole-etykieta';
    podpis.textContent = etykieta;

    element.append(podpis, wartosc);

    if (opis !== undefined) {
      const zdanie = document.createElement('p');
      zdanie.className = 'dn-pole-opis';
      zdanie.textContent = opis;
      element.append(zdanie);
    }

    return element;
  }

  function plakietka(tekst: string, wariant: 'informacja' | 'blad' | 'sygnal'): HTMLElement {
    const element = document.createElement('span');
    element.className = `dn-plakietka dn-plakietka--${wariant}`;
    element.textContent = tekst;
    return element;
  }

  function etykietaWartosci(isolated: boolean): string {
    return isolated ? 'Odrębna' : 'Współdzielona';
  }

  function wartoscPoczatkowaSelektu(stan: StanObecny, domyslnyIsolated: boolean): string {
    const isolated = stan.rodzaj === 'zapisane' ? stan.isolated : domyslnyIsolated;
    return isolated ? WARTOSC_ODREBNA : WARTOSC_WSPOLDZIELONA;
  }

  /**
   * Wiersz zmiany wartości: wybór nowej wartości (w pełni czynny) obok
   * przycisku zapisu, który woła naprawdę `isolation.context.set` i nanosi
   * odpowiedź albo odmowę — nigdy sukces przed odpowiedzią.
   */
  function zbudujWierszZmiany(punkt: PunktKontekstu, wartoscPoczatkowa: string): HTMLElement {
    const wybierz = wybor(`Nowa wartość — ${punkt.nazwa}`, OPCJE_WARTOSCI);
    wybierz.value = wartoscPoczatkowa;

    const poleZmiany = wiersz('Nowa wartość', wybierz, { klasa: 'dn-pole' });

    const zapisz = document.createElement('button');
    zapisz.type = 'button';
    zapisz.className = 'dn-btn dn-btn--zarys';
    zapisz.textContent = `Zapisz zmianę — ${punkt.nazwa}`;
    zapisz.addEventListener('click', () => void zapiszZmiane(punkt, wybierz.value));

    const przybornik = document.createElement('div');
    przybornik.className = 'dn-przybornik';
    przybornik.append(poleZmiany, zapisz);
    return przybornik;
  }

  async function zapiszZmiane(punkt: PunktKontekstu, wartosc: string): Promise<void> {
    const isolated = wartosc === WARTOSC_ODREBNA;

    const switches: IsolationSwitch[] = PUNKTY.map((p) => ({
      kind: p.kind,
      isolated: p.kind === punkt.kind ? isolated : (znane.get(p.kind) ?? true),
    }));

    const wynik = await posijKomende<{ switches: IsolationSwitch[] }>(
      kanal,
      Command.IsolationContextSet,
      {
        scope: zasieg.zasieg(),
        scopeId: zasieg.bytDoZadania(),
        layer: warstwa.warstwa(),
        switches,
      },
    );

    if (!wynik.udany || wynik.wynik === undefined) {
      tresc.potwierdzenie(
        `Zapis punktu „${punkt.nazwa}” nie powiódł się: ` +
          opisOdmowyBledu('Zapis punktu izolacji kontekstu', wynik.blad),
        false,
      );
      return;
    }

    tresc.potwierdzenie(
      `Punkt „${punkt.nazwa}” zapisany: ${etykietaWartosci(isolated)} — ${zasieg.opis()}, ` +
        `warstwa „${NAZWY_WARSTW[warstwa.warstwa()]}”.`,
      true,
    );
    void wczytaj();
  }

  /**
   * Objaśnienie kontekstowe [?] przełącznika — co ustawienie robi i jaki ma
   * wpływ na działanie aplikacji (rozdz. 3.3 Modelu konfiguracji). Zdanie
   * o skutku stoi obok opisu klucza, bo sama nazwa wartości („odrębna”)
   * nie mówi Operatorowi, co się po jej wybraniu zmieni w pracy.
   */
  function objasnienieKontekstowe(punkt: PunktKontekstu): string {
    return (
      `${punkt.objasnienie} Zapis obejmuje ${zasieg.opis()} i warstwę ` +
      `„${NAZWY_WARSTW[warstwa.warstwa()]}”; poziomy węższe od wskazanego zachowują własne zapisy, ` +
      'a poziomy szersze pozostają nietknięte. Stan wyjściowy platformy — odrębna — nie jest ' +
      'ograniczeniem możliwości, tylko uporządkowanym podziałem pracy; współdzielenie jest ' +
      'decyzją Operatora, nie ustawieniem narzuconym przez okno.'
    );
  }

  return { element: tresc.element, odswiez };
}

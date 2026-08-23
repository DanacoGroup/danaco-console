import { AgentVisibility, type Agent } from '../../../../shared/contract';
import { opisOdmowy } from '../../komponenty/odmowa';
import {
  poleLogiczne,
  poleTekstowe,
  poleWielowierszowe,
  poleWyboru,
  przycisk,
  utworzWierszOdpowiedzi,
} from '../../modele/kontrolki-formularza';
import { utworzPoziomyPamieci, type PoziomyPamieci } from './poziomy-pamieci';
import type { StanAgentow } from './stan-agentow';
import { czyZastepuje } from './warstwy-promptu';

/**
 * Formularz tożsamości eksperta — nazwa, imię własne, favikon, opis,
 * instrukcje, zasięg widoczności i poziomy pamięci.
 *
 * Zasięg i pamięć stoją tutaj, a nie w oknie osobnym, bo są komponentami
 * definicji zapisywanymi tą samą komendą co reszta tożsamości: jadą polami
 * `visibility` i `memoryLevels` żądania zakładającego i zmieniającego eksperta.
 * Osobne okno musiałoby wołać tę samą komendę drugi raz i zakładałoby drugą
 * wersję eksperta na każdą zmianę zasięgu.
 *
 * Moduł Agents jest kompozytorem: Operator nie konfiguruje tu konta ani API,
 * tylko nadaje surowemu modelowi tożsamość i zapisuje ją pod własną nazwą.
 *
 * Nazwa i imię to dwa pola. Nazwa jest nazwą bytu w bibliotece i po niej
 * ekspert odnajduje się w wykazie modułu. Imię własne jest tym, czym ekspert
 * przedstawia się w oknach roboczych całego produktu. Zlanie ich w jedno pole
 * odbierałoby Operatorowi możliwość nazwania stu agentów technicznie
 * (`referent-procesowy-v3`) i ludzko („Referent") jednocześnie.
 *
 * Odstępstwo od promptu globalnego stoi tutaj, a nie przy warstwach. Prompt
 * systemowy ustawia się globalnie w oknie konfiguracji na stronie głównej
 * i obowiązuje domyślnie; moduł Agents daje instrukcję dopisywaną do niego albo
 * jawne oznaczenie odstępstwa, po którym instrukcja eksperta staje się promptem
 * systemowym. Oznaczenie dotyczy instrukcji eksperta jako całości
 * (`Agent.mode`, jedno pole na eksperta), więc kontrolka stoi raz — przy
 * tożsamości. Warstwy (`warstwy-promptu.ts`) czytają ten stan i mówią o nim,
 * ale go nie ustawiają.
 *
 * Formularz nie buduje wybieraka emoji ani katalogu ikon: kontrakt niesie
 * `favicon` jako napis i nie ma komendy oddającej katalog znaków. Pole tekstowe
 * z podglądem obok mówi prawdę o tym, co pójdzie do rdzenia.
 */
/** Co robią instrukcje eksperta w stanie domyślnym — pole odznaczone. */
const SKUTEK_DOPISANIA =
  'Stan domyślny: instrukcje i warstwy tego eksperta DOPISUJĄ się do globalnego promptu ' +
  'systemowego ustawionego w oknie konfiguracji i ustawień na stronie głównej. Globalny ' +
  'obowiązuje pierwszy i zostaje w mocy.';

/** Co się stanie po oznaczeniu odstępstwa — pole zaznaczone. */
const SKUTEK_ZASTAPIENIA =
  'ODSTĘPSTWO OD USTAWIEŃ DOMYŚLNYCH: instrukcja tego eksperta STAJE SIĘ promptem ' +
  'systemowym. Globalny prompt z okna konfiguracji i ustawień przestaje obowiązywać ' +
  'dla okien obsługiwanych przez tego eksperta.';

export interface FormularzTozsamosci {
  element: HTMLElement;
  /** Nanosi eksperta czynnego na pola; `null` znaczy formularz zakładania. */
  ustaw(ekspert: Agent | null): void;
  /**
   * Dopisuje treść na końcu pola instrukcji, BEZ zapisu w rdzeniu.
   *
   * Rozdzielenie wklejenia od zapisu jest tu treścią, nie ostrożnością: rada
   * doradcy przeniesiona do instrukcji ma najpierw stanąć Operatorowi przed
   * oczami w polu, które sam potem zatwierdzi przyciskiem. Zapis wykonany
   * automatycznie zmieniłby tożsamość eksperta cudzym zdaniem, którego
   * Operator jeszcze nie przeczytał.
   */
  dopiszDoInstrukcji(tresc: string): void;
}

/** Zależności formularza. */
export interface OpcjeTozsamosci {
  /** Wywoływane po udanym zapisie — moduł odświeża wtedy bibliotekę. */
  naZapisie(ekspert: Agent): void;
}

export function utworzFormularzTozsamosci(
  stan: StanAgentow,
  opcje: OpcjeTozsamosci,
): FormularzTozsamosci {
  const nazwa = poleTekstowe({ etykieta: 'Nazwa eksperta', podpowiedz: 'np. Referent procesowy' });
  const imie = poleTekstowe({
    etykieta: 'Imię własne',
    podpowiedz: 'np. Referent',
    opis: 'Imię widziane w oknach roboczych; puste znaczy nazwę eksperta.',
  });
  const favikon = poleTekstowe({
    etykieta: 'Favikon',
    podpowiedz: 'znak, np. ⚙',
    opis: 'Znak graficzny eksperta w wykazie i w oknie. Podgląd stoi obok pola.',
  });
  const opis = poleTekstowe({ etykieta: 'Opis', podpowiedz: 'do czego ekspert służy' });
  const instrukcje = poleWielowierszowe(
    {
      etykieta: 'Instrukcje systemowe',
      opis:
        'Pole płaskie bytu agenta (Agent.systemPrompt). Warstwy promptu stoją niżej i są ' +
        'osobnym zapisem — to pole ich nie zastępuje ani nie powiela.',
    },
    6,
  );
  const czynny = poleLogiczne({
    etykieta: 'Ekspert czynny',
    opis:
      'Dotyczy bytu zapisanego. Przy zakładaniu pole nie jedzie do rdzenia — ' +
      'komenda zakładająca eksperta nie przyjmuje stanu czynności, a nowy ekspert ' +
      'powstaje czynny.',
  });
  czynny.kontrolka.checked = true;

  // Zasięg widoczności jest listą dwóch wartości kontraktu, a nie polem
  // logicznym: „projektowy” nie jest zaprzeczeniem „globalnego”, tylko innym
  // zasięgiem. Ekspert projektowy nie wskazuje projektu tutaj — przynależność
  // zapisuje przypisanie w module Workspace, a drugie miejsce zapisu byłoby
  // drugą prawdą o tym samym.
  const widocznosc = poleWyboru(
    {
      etykieta: 'Zasięg widoczności',
      opis:
        'Globalny znaczy „widoczny wszędzie” i jest stanem wyjściowym. Projektowy ' +
        'zawęża widoczność do projektów, do których ekspert został przypisany — ' +
        'przypisanie następuje w module Workspace, nie tutaj.',
    },
    [
      { wartosc: AgentVisibility.Global, etykieta: 'globalny — widoczny wszędzie' },
      { wartosc: AgentVisibility.Project, etykieta: 'projektowy — tylko w przypisanych projektach' },
    ],
  );

  const pamiec: PoziomyPamieci = utworzPoziomyPamieci();

  // Pole logiczne, a nie lista dwóch pozycji, bo te dwie wartości `Agent.mode`
  // nie są równorzędne: dopisanie jest stanem domyślnym całego produktu,
  // zastąpienie — odstąpieniem od ustawień globalnych ze strony głównej. Lista
  // postawiłaby obie na jednej półce; niezaznaczone pole mówi, że jedna z nich
  // jest normą.
  const zastepowanie = poleLogiczne({
    etykieta: 'Zastąp globalny prompt systemowy instrukcją tego eksperta',
  });

  // Ostrzeżenie stoi pod polem, na stałe i w obu stanach — nie chowa się za
  // najechaniem i nie jest zwijalne. Barwę stanu o donioślejszych skutkach
  // niesie `modele/tozsamosc.css`.
  const skutek = document.createElement('p');
  skutek.className = 'da-odstepstwo';
  skutek.setAttribute('role', 'note');

  zastepowanie.kontrolka.addEventListener('change', () => odswiezSkutek());

  function odswiezSkutek(): void {
    const zastepuje = zastepowanie.kontrolka.checked;
    skutek.textContent = zastepuje ? SKUTEK_ZASTAPIENIA : SKUTEK_DOPISANIA;
    skutek.dataset['zastapienie'] = String(zastepuje);
  }

  const zapisz = przycisk('Zapisz tożsamość', 'dn-btn dn-btn--sm dn-btn--atrament');
  const odpowiedz = utworzWierszOdpowiedzi();

  // Podgląd jest kontrolą: favikon bywa znakiem, którego pole tekstowe pokazuje
  // wąsko i w kroju pisma formularza, więc obok pola stoi ten sam znak
  // w rozmiarze, w jakim widać go w wykazie.
  const podglad = document.createElement('span');
  podglad.className = 'da-favikon__podglad';
  podglad.dataset['pusty'] = 'tak';
  podglad.textContent = '—';

  const pasFavikonu = document.createElement('div');
  pasFavikonu.className = 'da-favikon';
  pasFavikonu.append(favikon.element, podglad);

  favikon.kontrolka.addEventListener('input', () => odswiezPodglad());

  function odswiezPodglad(): void {
    const znak = favikon.kontrolka.value.trim();
    podglad.textContent = znak === '' ? '—' : znak;
    podglad.dataset['pusty'] = znak === '' ? 'tak' : 'nie';
  }

  const element = document.createElement('div');
  element.className = 'da-edytor__formularz';
  element.append(
    nazwa.element,
    imie.element,
    pasFavikonu,
    opis.element,
    instrukcje.element,
    zastepowanie.element,
    skutek,
    widocznosc.element,
    pamiec.element,
    czynny.element,
    zapisz,
    odpowiedz.element,
  );

  let edytowany: Agent | null = null;

  /** Wartość listy zasięgu przełożona na wyliczenie kontraktu. */
  function wybranaWidocznosc(): AgentVisibility {
    return widocznosc.kontrolka.value as AgentVisibility;
  }

  /**
   * Zapis rozstrzyga wybór, nie przycisk: brak eksperta czynnego znaczy
   * założenie nowego (`agent.create`), ekspert czynny znaczy zmianę
   * (`agent.update`). Dwóch osobnych przycisków nie ma, bo Operator nie ma
   * powodu rozstrzygać, której komendy użyć.
   */
  async function zapisanie(): Promise<void> {
    const wpisanaNazwa = nazwa.kontrolka.value.trim();
    if (wpisanaNazwa === '') {
      odpowiedz.pokaz('Ekspert wymaga nazwy — rdzeń odmówi zapisu bez niej.', false);
      return;
    }
    odpowiedz.pokaz('Zapis tożsamości w toku…', true);
    const wynik =
      edytowany === null
        ? await stan.zrodlo.utworz({
            nazwa: wpisanaNazwa,
            opis: opis.kontrolka.value,
            instrukcje: instrukcje.kontrolka.value,
            kanal: '',
            model: '',
            imie: imie.kontrolka.value,
            favikon: favikon.kontrolka.value,
            zastepuje: zastepowanie.kontrolka.checked,
            widocznosc: wybranaWidocznosc(),
            poziomyPamieci: pamiec.wybrane(),
          })
        : await stan.zrodlo.zmien(edytowany.id, {
            nazwa: wpisanaNazwa,
            opis: opis.kontrolka.value,
            instrukcje: instrukcje.kontrolka.value,
            czynny: czynny.kontrolka.checked,
            imie: imie.kontrolka.value,
            favikon: favikon.kontrolka.value,
            zastepuje: zastepowanie.kontrolka.checked,
            widocznosc: wybranaWidocznosc(),
            poziomyPamieci: pamiec.wybrane(),
          });

    if (!wynik.udany || wynik.wynik === undefined) {
      odpowiedz.pokaz(opisOdmowy('Zapis tożsamości', wynik.blad?.code, wynik.blad?.message), false);
      return;
    }
    const zapisany = wynik.wynik.agent;
    stan.wchlon(zapisany);
    stan.wybierz(zapisany.id);
    opcje.naZapisie(zapisany);
    odpowiedz.pokaz(
      `Tożsamość zapisana — ekspert „${zapisany.displayName ?? zapisany.name}", ` +
        `wersja ${zapisany.version ?? 1}.`,
      true,
    );
  }

  zapisz.addEventListener('click', () => void zapisanie());

  return {
    element,

    ustaw(ekspert) {
      edytowany = ekspert;
      nazwa.kontrolka.value = ekspert?.name ?? '';
      imie.kontrolka.value = ekspert?.displayName ?? '';
      favikon.kontrolka.value = ekspert?.favicon ?? '';
      opis.kontrolka.value = ekspert?.description ?? '';
      instrukcje.kontrolka.value = ekspert?.systemPrompt ?? '';
      czynny.kontrolka.checked = ekspert?.enabled ?? true;
      // Ekspert bez `mode` jest ekspertem dopisującym — pominięte pole znaczy
      // w kontrakcie `DOLACZ`, więc formularz nowego eksperta rusza odznaczony.
      zastepowanie.kontrolka.checked = czyZastepuje(ekspert);
      odswiezSkutek();
      widocznosc.kontrolka.value = ekspert?.visibility ?? AgentVisibility.Global;
      // Poziomy pamięci są w kontrakcie polem wymaganym eksperta, więc przy
      // ekspercie zapisanym idą wprost z rdzenia; `null` znaczy formularz
      // zakładania i wtedy grupa rusza od stanu wyjściowego.
      pamiec.ustaw(ekspert === null ? null : ekspert.memoryLevels);
      zapisz.textContent = ekspert === null ? 'Załóż eksperta' : 'Zapisz tożsamość';
      odswiezPodglad();
    },

    dopiszDoInstrukcji(tresc) {
      const biezaca = instrukcje.kontrolka.value;
      instrukcje.kontrolka.value = biezaca === '' ? tresc : `${biezaca}\n\n${tresc}`;
      instrukcje.kontrolka.focus();
      odpowiedz.pokaz(
        'Rada doradcy wklejona do instrukcji wraz z nagłówkiem prowenancji. ' +
          'Rdzeń jej jeszcze nie zna — zapis idzie przyciskiem powyżej.',
        true,
      );
    },
  };
}

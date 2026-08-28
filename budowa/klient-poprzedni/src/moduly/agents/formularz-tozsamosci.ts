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
 * Formularz tożsamości eksperta zbiera nazwę, imię własne, favikon, opis,
 * instrukcje, zasięg widoczności oraz poziomy pamięci zapisywane jedną
 * komendą kontraktu.
 */
const SKUTEK_DOPISANIA =
  'Stan domyślny: instrukcje i warstwy tego eksperta DOPISUJĄ się do globalnego promptu ' +
  'systemowego ustawionego w oknie konfiguracji i ustawień na stronie głównej. Globalny ' +
  'obowiązuje pierwszy i zostaje w mocy.';

const SKUTEK_ZASTAPIENIA =
  'ODSTĘPSTWO OD USTAWIEŃ DOMYŚLNYCH: instrukcja tego eksperta STAJE SIĘ promptem ' +
  'systemowym. Globalny prompt z okna konfiguracji i ustawień przestaje obowiązywać ' +
  'dla okien obsługiwanych przez tego eksperta.';

export interface FormularzTozsamosci {
  element: HTMLElement;
  /** Nanosi eksperta czynnego na pola; `null` znaczy formularz zakładania. */
  ustaw(ekspert: Agent | null): void;
  /** Dopisuje treść na końcu pola instrukcji bez zapisu — zapis zatwierdza Operator przyciskiem. */
  dopiszDoInstrukcji(tresc: string): void;
}

/** Zależności formularza przekazywane z zewnątrz, wywoływane po zdarzeniach formularza, którymi moduł nadrzędny steruje odświeżeniem biblioteki. */
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

  // Zasięg widoczności jest listą dwóch wartości kontraktu, nie polem logicznym — to odrębne zasięgi.
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

  // Pole logiczne, nie lista pozycji, bo dopisanie jest stanem domyślnym, a zastąpienie — odstąpieniem.
  const zastepowanie = poleLogiczne({
    etykieta: 'Zastąp globalny prompt systemowy instrukcją tego eksperta',
  });

  // Ostrzeżenie stoi pod polem na stałe w obu stanach, bez ukrywania go za najechaniem kursora.
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

  // Podgląd pokazuje ten sam znak favikonu w rozmiarze widocznym w wykazie, obok pola tekstowego.
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

  /** Zapis rozstrzyga wybór: brak eksperta czynnego zakłada nowego, ekspert czynny zmienia istniejącego. */
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
      // Ekspert bez trybu jest ekspertem dopisującym — pominięte pole w kontrakcie znaczy dopisanie.
      zastepowanie.kontrolka.checked = czyZastepuje(ekspert);
      odswiezSkutek();
      widocznosc.kontrolka.value = ekspert?.visibility ?? AgentVisibility.Global;
      // Poziomy pamięci eksperta zapisanego idą z rdzenia; puste oznacza formularz zakładania nowego.
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

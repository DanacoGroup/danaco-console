import { Command } from '../../../../shared/contract';
import { opisOdmowyBledu } from '../../komponenty/odmowa';
import { przycisk, utworzWierszOdpowiedzi } from '../../modele/kontrolki-formularza';
import { OKNO_PROMPT_BUILDER } from './etykiety-designu';
import { utworzHistoriePromptow, type HistoriaPromptow } from './historia-promptow';
import { utworzModalKreatora, type ModalKreatora } from './modal-kreatora';
import { utworzPasekPostepu, type PasekPostepu } from './pasek-postepu';
import { utworzPolaPromptu, type PolaPromptu } from './pola-promptu';
import { utworzSzablonyPromptu, type SzablonyPromptu } from './szablony-promptu';
import { polecenieZOdmowy, zdanieOPoleceniu } from './polecenie-z-odmowy';
import { naglowekOkna, utworzStanOkna, type StanOkna } from './stan-okna';
import { powodZTorem } from './tor-komendy';
import type { StanDesignu } from './stan-designu';

/**
 * Prompt Builder — okno **kreator** modułu Design (kod katalogu rdzenia
 * `prompt-builder`): złożenie promptu strukturalnego i wydanie zlecenia
 * generowania.
 *
 * Pola strukturalne mieszkają w modalu, bo klasa okna „Kreator" żąda nośnika
 * modalnego w stanie zamkniętym; sekcja okna zostaje widoczna zawsze, bo jest
 * jednym z okien operacyjnych modułu. Przycisk `⚡ Generuj` stoi w obu miejscach
 * i prowadzi do tej samej czynności.
 *
 * Generowanie kończy się jedną z dwóch dróg i okno obsługuje obie.
 * `adapter_modul_design_generowanie.go` wysyła polecenie kanałem obrazowym,
 * odbiera fragment `image`, odkłada bajty w magazynie pod sumą sha256, mierzy
 * format i wymiary z nagłówka utrwalonego pliku i zakłada wiersz zasobu z `uri`
 * — odpowiedź udana niesie zasoby z treścią. Odmowa przychodzi przy braku
 * kanału obrazowego w rejestrze, kanale nieczynnym, kanale tekstowym, braku
 * poświadczenia albo odpowiedzi bez obrazu; niesie wtedy
 * `error.details.polecenie` z gotową treścią polecenia, więc okno pokazuje
 * odmowę wraz z oddanym tekstem — złożenie promptu odbyło się także wtedy.
 *
 * Odpowiedź udanego generowania ma w kontrakcie pole `processId`, którego rdzeń
 * nie wypełnia, a zdarzenia `progress.changed` nie przychodzą. Bez `processId`
 * pasek postępu jest więc wyciszany, zamiast stać na „zlecenie przyjęte" po
 * pracy już zakończonej.
 */
export interface OknoPromptBuilder {
  element: HTMLElement;
  odswiez(): void;
  /** Przyjmuje zdarzenie postępu z magistrali modułu. */
  przyjmijPostep(...argumenty: Parameters<PasekPostepu['przyjmij']>): void;
}

export function utworzOknoPromptBuilder(stan: StanDesignu): OknoPromptBuilder {
  const okno: StanOkna = utworzStanOkna();
  const pola: PolaPromptu = utworzPolaPromptu();
  const postep: PasekPostepu = utworzPasekPostepu();
  const modal: ModalKreatora = utworzModalKreatora('Prompt Builder — prompt strukturalny', '⚡ Generuj');
  const historia: HistoriaPromptow = utworzHistoriePromptow((prompt) => {
    pola.nanies(prompt);
    modal.otworz();
  });

  // Szablony i historia promptów trwałej: obie żyły dotąd w oknie do zamknięcia
  // karty przeglądarki. Panel oddaje prompt z powrotem w pola kreatora, więc
  // szablon da się użyć, a nie tylko obejrzeć.
  const szablony: SzablonyPromptu = utworzSzablonyPromptu(stan, {
    prompt: () => pola.prompt(),
    naSzablon: (prompt) => {
      pola.nanies(prompt);
      modal.otworz();
    },
  });

  const licznik = document.createElement('p');
  licznik.className = 'dn-pole-opis md-prompt__licznik';

  const otworz = przycisk('Otwórz kreator promptu', 'dn-btn dn-btn--sm dn-btn--zarys');
  const generuj = przycisk('⚡ Generuj', 'dn-btn dn-btn--sm dn-btn--sygnal');
  const odpowiedz = utworzWierszOdpowiedzi();

  modal.cialo.append(pola.element, licznik);
  okno.tresc.append(
    otworz,
    generuj,
    odpowiedz.element,
    postep.element,
    historia.element,
    szablony.element,
    modal.element,
  );

  const element = document.createElement('section');
  element.className = 'md-okno md-okno--kreator';
  element.dataset['okno'] = OKNO_PROMPT_BUILDER.kod;
  element.append(naglowekOkna(OKNO_PROMPT_BUILDER.nazwa, OKNO_PROMPT_BUILDER.rola), okno.element);

  /** Komunikat blokowy modala; zostaje po zamknięciu, aż do kolejnej próby. */
  function powiedz(tresc: string, powodzenie: boolean): void {
    modal.komunikat.hidden = tresc === '';
    modal.komunikat.textContent = tresc;
    modal.komunikat.dataset['powodzenie'] = String(powodzenie);
    odpowiedz.pokaz(tresc, powodzenie);
  }

  async function generujZasob(): Promise<void> {
    const prompt = pola.prompt();
    // Próba z brakami: przycisk był czynny od otwarcia, więc naciśnięcie ma
    // ujawnić brak, a kreator ma zostać otwarty.
    if (prompt.subject === '') {
      powiedz('Temat jest jedynym polem wymaganym promptu — bez niego rdzeń odmówi zlecenia.', false);
      modal.otworz();
      return;
    }
    if (stan.idOkna() === '') {
      powiedz(`Komenda ${Command.DesignAssetGenerate} wymaga okna modułu. ${stan.opisOkna()}`, false);
      return;
    }
    okno.ladowanie('Zlecenie generowania w toku…');
    powiedz('Zlecenie generowania wysłane — czekam na odpowiedź rdzenia.', true);
    postep.oczekuj('', stan.idOkna());
    // Czuwanie pilnuje, żeby okno nie kręciło wskaźnika po zerwanym gnieździe.
    // Cisza kanału nie jest odmową generowania: rdzeń mógł zlecenie odebrać
    // i wykonać, więc pasek postępu milknie, a zdanie mówi o skutku nieznanym,
    // nie o niepowodzeniu.
    const wynik = await stan.czuwanie.prowadz(
      'generowanie zasobu',
      stan.zrodlo.generuj({
        idOkna: stan.idOkna(),
        prompt,
        idReferencji: pola.idReferencji(),
        // Wskazany kanał jedzie polem, po którym rdzeń wybiera adapter, a nie
        // polem opisowym.
        idKanalu: pola.idKanalu(),
        rodzaj: pola.rodzaj(),
      }),
      {
        wToku: (zdanie) => okno.ladowanie(zdanie),
        cisza: (zdanie) => {
          okno.blad(zdanie);
          powiedz(zdanie, false);
          postep.wycisz();
        },
        powrot: (zdanie) => {
          okno.blad(zdanie);
          powiedz(zdanie, false);
        },
        spozniona: (zdanie) => powiedz(zdanie, false),
      },
    );
    // `null` znaczy „bez rozstrzygnięcia" — zdanie stoi już w oknie i w modalu,
    // a rdzeń mógł zlecenie wykonać, więc odmowy się tu nie dopisuje.
    if (wynik === null) return;
    if (!wynik.udany || wynik.wynik === undefined) {
      const powod = opisOdmowyBledu('Generowanie zasobu', wynik.blad);
      okno.blad(powodZTorem(powod, Command.DesignAssetGenerate));
      // Odmowa braku silnika niesie w `details` gotową treść polecenia. Prompt
      // trafia wtedy do historii, bo złożenie promptu się odbyło — inaczej niż
      // przy odmowie wskazania (bez okna, bez tematu), gdzie składać nie było co.
      const polecenie = polecenieZOdmowy(wynik.blad);
      if (polecenie !== '') {
        historia.dopisz(prompt);
        powiedz(`${powod}\n\n${zdanieOPoleceniu(polecenie)}`, false);
      } else {
        powiedz(powod, false);
      }
      postep.wycisz();
      return;
    }
    historia.dopisz(prompt);
    for (const zasob of wynik.wynik.assets) stan.wchlon(zasob);
    // Bez identyfikatora procesu nie ma czego śledzić, więc pasek milknie
    // zamiast zostawać na „czekam" po zakończonej turze.
    const idProcesu = (wynik.wynik.processId ?? '').trim();
    if (idProcesu === '') postep.wycisz();
    else postep.oczekuj(idProcesu, stan.idOkna());
    okno.gotowe();
    powiedz(
      `Rdzeń oddał ${wynik.wynik.assets.length} zasób/zasoby do Assets Panel.`,
      true,
    );
  }

  otworz.addEventListener('click', () => modal.otworz());
  for (const kontrolka of [generuj, modal.glowny]) {
    kontrolka.addEventListener('click', () => void generujZasob());
  }

  function odswiez(): void {
    pola.odswiez(stan.silniki(), stan.zasoby(), historia.style());
    licznik.textContent = `Długość treści promptu: ${pola.dlugosc()} znaków.`;
    if (okno.faza() === 'ladowanie' || okno.faza() === 'blad') return;
    okno.puste(
      modal.otwarty()
        ? 'Kreator otwarty — wypełnij pola strukturalne i naciśnij „⚡ Generuj".'
        : 'Kreator zamknięty. Otwórz go, aby złożyć prompt strukturalny.',
    );
  }

  odswiez();

  return {
    element,
    odswiez,
    przyjmijPostep: (tresc) => postep.przyjmij(tresc),
  };
}

import { poleWielowierszowe, przycisk, utworzWierszOdpowiedzi } from '../../modele/kontrolki-formularza';
import { AKCJE_WORKSPACE } from './akcje-okien';
import {
  ustawStanZakresu,
  ustawZakres,
  wykonajAkcjeZakresu,
  type KontekstZakresu,
} from './czynnosci-zakresu';
import { zDymkiem } from './dymek-badania';
import { KODY_OKIEN } from './kody-okien';
import { utworzPanelPostepu } from './panel-postepu';
import { utworzPasekEtapow } from './pasek-etapow';
import { utworzRameBadania } from './rama-badania';
import type { StanBadania } from './stan-badania';
import { utworzStanOknaBadania } from './stan-okna-badania';

/**
 * Research Workspace — okno wiodące modułu: zakres badania i wejście do
 * czterech pozostałych okien.
 *
 * Zakres badania zapisuje komenda `research.workspace.set` (zakres i etapy),
 * której uchwyt stoi w rdzeniu — `adapter_modul_badania_uchwyty.go`,
 * `CommandResearchWorkspaceSet`. Zapis kończy się odpowiedzią rdzenia albo jego
 * odmową, nigdy ciszą.
 *
 * Nawigacja do pozostałych okien modułu jest przeniesieniem ogniska wewnątrz
 * przestrzeni modułu: okna stoją obok siebie, nie w osobnych trasach, więc nie
 * woła rdzenia.
 */
export interface OknoResearchWorkspace {
  element: HTMLElement;
  odswiez(): void;
}

/**
 * Wykaz okien, do których to okno przenosi ognisko — w porządku pracy badawczej
 * z opracowania (rozdz. 4.1), nie w porządku alfabetycznym.
 */
const PRZEJSCIA: readonly { kod: string; nazwa: string }[] = [
  { kod: KODY_OKIEN.odkrywanie, nazwa: 'Discovery Panel' },
  { kod: KODY_OKIEN.zrodla, nazwa: 'Sources Manager' },
  { kod: KODY_OKIEN.lektura, nazwa: 'Reading View' },
  { kod: KODY_OKIEN.ustalenia, nazwa: 'Findings Panel' },
  { kod: KODY_OKIEN.raport, nazwa: 'Report Builder' },
  { kod: KODY_OKIEN.eksport, nazwa: 'Export Panel' },
];

export function utworzOknoResearchWorkspace(
  stan: StanBadania,
  przejdz: (kodOkna: string) => void,
): OknoResearchWorkspace {
  const poleZakresu = poleWielowierszowe(
    { etykieta: 'Zakres badania', podpowiedz: 'co i po co badamy' },
    3,
  );
  const poleEtapow = poleWielowierszowe(
    { etykieta: 'Etapy badania', podpowiedz: 'jeden etap w wierszu' },
    4,
  );

  const kontekst: KontekstZakresu = {
    stan,
    okno: utworzStanOknaBadania(),
    odpowiedz: utworzWierszOdpowiedzi(),
    zakres: poleZakresu.kontrolka,
    etapy: poleEtapow.kontrolka,
    przejdz,
    odswiez: () => odswiez(),
  };

  const zapisz = przycisk('Zapisz zakres badania', 'dn-btn dn-btn--sm dn-btn--atrament');
  zapisz.addEventListener('click', () => void ustawZakres(kontekst));

  const nawigacja = document.createElement('ul');
  nawigacja.className = 'mr-przejscia';
  nawigacja.append(...PRZEJSCIA.map((przejscie) => pozycjaPrzejscia(przejscie, przejdz)));

  // Wskazanie etapu prowadzi wzrok do okna właściwego temu krokowi pracy.
  // Opracowanie chce przewinięcia „do materiału właściwego etapowi". Przypisanie
  // źródła do etapu ma już pole w żądaniu katalogowania, ale samo źródło
  // oddawane przez rdzeń go nie niesie, więc okno nie ma po czym zawężać
  // i mówi to wprost, zamiast udawać filtr.
  const etapy = utworzPasekEtapow((etap) => {
    kontekst.odpowiedz.pokaz(
      `Wskazano etap „${etap}". Przypisanie źródła do etapu da się wysłać, ale źródło oddawane ` +
        'przez rdzeń go nie niesie, więc wskazanie nie zawęża materiału — porządkuje pracę ' +
        'w tym widoku.',
      true,
    );
  });

  const postep = utworzPanelPostepu(stan, przejdz);

  kontekst.okno.tresc.append(
    zDymkiem(poleZakresu.element, 'Zakres jedzie w polu scope komendy research.workspace.set.'),
    zDymkiem(poleEtapow.element, 'Etapy jadą w polu stages tej samej komendy; puste wiersze są pomijane.'),
    zapisz,
    kontekst.odpowiedz.element,
    etapy.element,
    postep.element,
    nawigacja,
  );

  const rama = utworzRameBadania(
    KODY_OKIEN.workspace,
    'Research Workspace',
    'wiodące',
    AKCJE_WORKSPACE,
    (akcja) => void wykonajAkcjeZakresu(kontekst, akcja),
  );
  rama.cialo.append(kontekst.okno.element);

  function odswiez(): void {
    etapy.odswiez(stan.etapy());
    postep.odswiez();
    ustawStanZakresu(kontekst);
  }

  odswiez();
  return { element: rama.element, odswiez };
}

/** Pozycja nawigacji do okna modułu; w pełni klikalna zawsze. */
function pozycjaPrzejscia(
  przejscie: { kod: string; nazwa: string },
  przejdz: (kodOkna: string) => void,
): HTMLElement {
  const element = document.createElement('li');
  element.className = 'mr-przejscia__wiersz';

  const kontrolka = przycisk(przejscie.nazwa, 'dn-btn dn-btn--sm dn-btn--duch');
  kontrolka.dataset['przejscie'] = przejscie.kod;
  kontrolka.addEventListener('click', () => przejdz(przejscie.kod));

  element.append(kontrolka);
  return element;
}

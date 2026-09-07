// Panel czynności w oknie Assistant: czynności okna komunikacji, przestawienie
// samego okna i przekazanie pracy innemu oknu.
import { Command, ExecutionEnv, PermissionMode } from '../../../shared/contract.ts';
import type { Kanal } from '../protokol/kanal.ts';
import { wywolaj } from '../protokol/wywolanie.ts';
import { dolozCzynnosciPanelu, zapytajWSzufladzie } from './czynnosci-okna.ts';
import { oglos } from './ogloszenie.ts';

const NAGLOWEK = 'Okno komunikacji';
const PANEL = 'panel-actions';

const SPOSOBY_ZGODY: ReadonlyArray<readonly [string, string]> = [
  [PermissionMode.Manual, 'Pytaj o każdą zmianę'],
  [PermissionMode.AcceptEdits, 'Przyjmuj poprawki plików'],
  [PermissionMode.Plan, 'Tylko plan, bez zmian'],
  [PermissionMode.Auto, 'Samodzielnie'],
  [PermissionMode.DontAsk, 'Bez pytania'],
  [PermissionMode.BypassPermissions, 'Z pominięciem zgód'],
];

const SRODOWISKA: ReadonlyArray<readonly [string, string]> = [
  [ExecutionEnv.Local, 'Maszyna Operatora'],
  [ExecutionEnv.Core, 'Rdzeń'],
  [ExecutionEnv.Remote, 'Maszyna zdalna'],
];

interface Otoczenie {
  kanal: Kanal;
  korzen: Element;
  idOkna: () => string;
}

export function zwiazOknoAsystenta(
  kanal: Kanal,
  korzen: Element,
  idOkna: () => string,
  dolacz: (zdejmij: () => void) => void,
  przy: AddEventListenerOptions,
): void {
  const otoczenie: Otoczenie = { kanal, korzen, idOkna };
  const zdejmij = dolozCzynnosciPanelu(korzen, PANEL, 'Czynności okna', [
    {
      naglowek: 'Czynności',
      pozycje: [
        { kod: 'wykaz', nazwa: 'Wykaz czynności okna' },
        { kod: 'wykonaj', nazwa: 'Wykonaj czynność…' },
      ],
    },
    {
      naglowek: 'Samo okno',
      pozycje: [
        { kod: 'przestaw', nazwa: 'Przestaw okno…' },
        { kod: 'przekaz', nazwa: 'Przekaż pracę innemu oknu…' },
      ],
    },
  ], (kod) => {
    void wykonaj(otoczenie, kod);
  }, przy);
  if (zdejmij !== null) dolacz(zdejmij);
}

function wypelnij(korzen: Element, wiersze: string[]): void {
  const cialo = korzen.querySelector(`#${PANEL} .sta-okno-tresc`);
  if (cialo === null) return;
  cialo.replaceChildren(...wiersze.map((tresc) => {
    const wiersz = cialo.ownerDocument.createElement('div');
    wiersz.className = 'dn-wykaz-modulu-poz';
    wiersz.textContent = tresc;
    return wiersz;
  }));
}

async function wykonaj(otoczenie: Otoczenie, kod: string): Promise<void> {
  if (otoczenie.idOkna() === '') {
    oglos(NAGLOWEK, 'To okno nie stoi jeszcze przy sesji.', 'ostrzezenie');
    return;
  }
  if (kod === 'wykaz') return wykazCzynnosci(otoczenie);
  if (kod === 'wykonaj') return wykonajCzynnosc(otoczenie);
  if (kod === 'przestaw') return przestawOkno(otoczenie);
  if (kod === 'przekaz') return przekazPrace(otoczenie);
  /* Czynność spoza obsłużonych odmawia zamiast milczeć. */
  oglos(NAGLOWEK, `Czynność „${kod}" nie jest prowadzona przez to okno.`, 'ostrzezenie');
}

async function czynnosciOkna(kanal: Kanal): Promise<ReadonlyArray<readonly [string, string]>> {
  const wykaz = await wywolaj(kanal, Command.ActionList, { enabledOnly: true });
  return (wykaz.wynik?.actions ?? []).map((czynnosc) => [czynnosc.id, czynnosc.name] as const);
}

async function wykazCzynnosci(otoczenie: Otoczenie): Promise<void> {
  const wynik = await wywolaj(otoczenie.kanal, Command.ActionList, { enabledOnly: false });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił wykazu czynności.', 'ostrzezenie');
    return;
  }
  const czynnosci = wynik.wynik.actions;
  wypelnij(otoczenie.korzen, czynnosci.length === 0
    ? ['Rdzeń nie podaje żadnej czynności okna.']
    : czynnosci.map((czynnosc) => `${czynnosc.name} · ${czynnosc.command}`
      + (czynnosc.enabled ? '' : ' · wyłączona')));
  oglos(NAGLOWEK, `Czynności w wykazie: ${String(czynnosci.length)}.`);
}

/* Parametry czynności rdzeń bierze jako zapis JSON, bo każda czynność ma swój
   kształt; puste pole znaczy czynność bez parametrów. */
async function wykonajCzynnosc(otoczenie: Otoczenie): Promise<void> {
  const wybor = await czynnosciOkna(otoczenie.kanal);
  if (wybor.length === 0) {
    oglos(NAGLOWEK, 'Rdzeń nie podaje żadnej czynności czynnej.', 'ostrzezenie');
    return;
  }
  const wartosci = await zapytajWSzufladzie(otoczenie.korzen, PANEL, {
    tytul: 'Wykonanie czynności okna',
    pola: [
      { klucz: 'czynnosc', etykieta: 'Czynność', wybor },
      { klucz: 'parametry', etykieta: 'Parametry (zapis JSON)', obszerne: true },
    ],
    wykonanie: 'Wykonaj',
  });
  if (wartosci === null) return;
  const zapis = (wartosci.parametry ?? '').trim();
  let parametry: unknown;
  if (zapis !== '') {
    try {
      parametry = JSON.parse(zapis);
    } catch {
      oglos(NAGLOWEK, 'Parametry nie są poprawnym zapisem JSON.', 'ostrzezenie');
      return;
    }
  }
  /* Czynność sięga do rdzenia i potrafi trwać, więc okno mówi o tym od razu —
     inaczej naciśnięcie wygląda na nieprzyjęte. */
  oglos(NAGLOWEK, 'Zlecam czynność rdzeniowi…');
  const wynik = await wywolaj(otoczenie.kanal, Command.WindowAction, {
    windowId: otoczenie.idOkna(),
    actionId: wartosci.czynnosc ?? '',
    ...(parametry === undefined ? {} : { parameters: parametry }),
  });
  if (!wynik.udany) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił wykonania czynności.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, 'Czynność okna wykonana.');
}

async function przestawOkno(otoczenie: Otoczenie): Promise<void> {
  const wartosci = await zapytajWSzufladzie(otoczenie.korzen, PANEL, {
    tytul: 'Przestawienie okna',
    opis: 'Puste pole zostawia nastawę okna nietkniętą.',
    pola: [
      { klucz: 'nazwa', etykieta: 'Nazwa okna' },
      { klucz: 'katalogi', etykieta: 'Katalogi robocze', obszerne: true },
      { klucz: 'zgoda', etykieta: 'Sposób pytania o zgodę', wybor: SPOSOBY_ZGODY },
      { klucz: 'srodowisko', etykieta: 'Miejsce wykonania', wybor: SRODOWISKA },
    ],
    wykonanie: 'Przestaw okno',
  });
  if (wartosci === null) return;
  const katalogi = (wartosci.katalogi ?? '').split('\n')
    .map((wiersz) => wiersz.trim())
    .filter((wiersz) => wiersz !== '');
  const wynik = await wywolaj(otoczenie.kanal, Command.WindowUpdate, {
    windowId: otoczenie.idOkna(),
    ...(wartosci.nazwa === '' ? {} : { title: wartosci.nazwa }),
    ...(katalogi.length === 0 ? {} : { workingDirs: katalogi }),
    permissionMode: (wartosci.zgoda ?? PermissionMode.Manual) as PermissionMode,
    executionEnv: (wartosci.srodowisko ?? ExecutionEnv.Local) as ExecutionEnv,
  });
  if (!wynik.udany) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił przestawienia okna.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, 'Okno przestawione.');
}

/* Przekazanie niesie polecenie dla okna przejmującego: bez niego okno dostaje
   pracę, o której nic nie wie. */
async function przekazPrace(otoczenie: Otoczenie): Promise<void> {
  const stan = await wywolaj(otoczenie.kanal, Command.WindowStateGet, {
    windowId: otoczenie.idOkna(),
  });
  const idSesji = stan.wynik?.window.sessionId ?? '';
  if (idSesji === '') {
    oglos(NAGLOWEK, 'Rdzeń nie podaje sesji tego okna.', 'ostrzezenie');
    return;
  }
  const okna = await wywolaj(otoczenie.kanal, Command.WindowList, { sessionId: idSesji });
  const wybor = (okna.wynik?.windows ?? [])
    .filter((okno) => okno.id !== otoczenie.idOkna())
    .map((okno) => [okno.id, okno.title ?? okno.id] as const);
  if (wybor.length === 0) {
    oglos(NAGLOWEK, 'Ta sesja nie ma drugiego okna, któremu można przekazać pracę.',
      'ostrzezenie');
    return;
  }
  const wartosci = await zapytajWSzufladzie(otoczenie.korzen, PANEL, {
    tytul: 'Przekazanie pracy',
    pola: [
      { klucz: 'okno', etykieta: 'Okno przejmujące', wybor },
      { klucz: 'polecenie', etykieta: 'Polecenie dla okna', obszerne: true, wymagane: true },
    ],
    wykonanie: 'Przekaż pracę',
  });
  if (wartosci === null) return;
  const wynik = await wywolaj(otoczenie.kanal, Command.WindowHandoff, {
    sessionId: idSesji,
    fromWindowId: otoczenie.idOkna(),
    toWindowId: wartosci.okno ?? '',
    instruction: wartosci.polecenie ?? '',
  });
  if (!wynik.udany) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił przekazania pracy.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, 'Praca przekazana wskazanemu oknu.');
}

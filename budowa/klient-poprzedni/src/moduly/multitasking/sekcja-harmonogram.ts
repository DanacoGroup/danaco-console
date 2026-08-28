import type { AutomationSchedule } from '../../../../shared/contract';
import {
  poleLogiczne,
  poleTekstowe,
  pozycjaWykazu,
  przyciskAkcji as przycisk,
  wykaz,
} from '../../modele/kontrolki-formularza';
import type { ZrodloSekcjiPaneli } from '../../powloka/zrodlo-sekcji-paneli';
import { PRZEDROSTEK, wierszOpisu } from './kontrolki';
import { utworzPowierzchnieSekcji, zdaniePuste, zObjasnieniem } from './powierzchnia-sekcji';
import { utworzWiezAutomations } from './wiez-automations';
import type { ZrodloNadzoru } from './zrodlo-nadzoru';

/**
 * Sekcja harmonogramu i automatyki zakłada regułę czasową pracy ciągłej na
 * wskazanym układzie automatyki, bo bez wskazania nie ma czego uruchamiać
 * cyklicznie.
 */
export interface SekcjaHarmonogramu {
  element: HTMLElement;
  odswiez(): void;
  ustawOkno(idOkna: string): void;
}

export interface OpcjeSekcjiHarmonogramu {
  nadzor: ZrodloNadzoru;
  sekcje: ZrodloSekcjiPaneli;
}

/**
 * Kształt komendy, której kontrakt nie ma: zdjęcie jednego wyzwalacza, bo
 * żadna komenda nie zdejmuje pojedynczego wyzwalacza osobno od zapisu całej
 * reguły.
 */
const BRAK_WYZWALACZY =
  'schedule.trigger.remove { scheduleId, triggerId } → { schedule: AutomationSchedule } — kontrakt komendy nie ma; wyzwalacz zdejmuje się dziś wyłącznie przez ponowny zapis całej reguły.';

export function utworzSekcjeHarmonogramu(
  opcje: OpcjeSekcjiHarmonogramu,
): SekcjaHarmonogramu {
  const { nadzor, sekcje } = opcje;

  const cykl = poleTekstowe({
    etykieta: 'Cykliczność (notacja cron)',
    podpowiedz: '0 * * * *',
  });
  const strefa = poleTekstowe({
    etykieta: 'Strefa czasowa',
    podpowiedz: 'Europe/Warsaw',
  });
  const obowiazuje = poleLogiczne({ etykieta: 'Harmonogram obowiązuje' });
  const zapisz = przycisk('Zapisz regułę czasową', 'dn-btn dn-btn--sm');

  const formularz = document.createElement('div');
  formularz.className = 'dm-orkiestracja__pasek';
  formularz.append(
    zObjasnieniem(
      cykl.element,
      'Reguła czasowa w notacji cron: co ile i o której rdzeń ma uruchamiać wskazany układ. To ona zamienia zespół pracujący na żądanie w pracę ciągłą 24/7 — bez niej układ rusza wyłącznie ręcznie.',
    ),
    zObjasnieniem(
      strefa.element,
      'Strefa, w której czytana jest notacja cron. Pominięta znaczy strefę rdzenia — na instalacji w innej strefie ta sama reguła uruchomi pracę o innej godzinie.',
    ),
    zObjasnieniem(
      obowiazuje.element,
      'Włącza i wyłącza regułę bez jej kasowania. Wyłączona zostaje zapisana w rdzeniu i daje się przywrócić jednym naciśnięciem — praca ciągła jest decyzją odwracalną.',
    ),
    zapisz,
  );

  const lista = wykaz('Harmonogramy układu', 'dm-wykaz');
  const braki = wykaz('Czego kontrakt nie niesie', 'dm-wykaz');
  braki.replaceChildren(
    pozycjaWykazu('Zdjęcie jednego wyzwalacza', BRAK_WYZWALACZY, PRZEDROSTEK).element,
  );

  const wiez = utworzWiezAutomations({
    okno: 'Scheduler (oraz Execution Monitor)',
    rola: 'reguła czasowa pracy ciągłej i podgląd jej przebiegów.',
  });

  const powierzchnia = utworzPowierzchnieSekcji({
    klucz: 'harmonogram',
    tytul: 'Harmonogram i automatyki',
    zakres:
      'Reguła czasowa pracy ciągłej 24/7/365 i wpięcie automatyk. Integracja z modułem Automations jest jawną, odwracalną decyzją Operatora.',
    sekcje,
    podsekcje: [
      { id: 'wiez', tytul: 'Powiązanie z modułem Automations', tresc: wiez.element },
      { id: 'regula', tytul: 'Reguła czasowa', tresc: formularz },
      { id: 'harmonogramy', tytul: 'Harmonogramy układu', tresc: lista },
      { id: 'braki', tytul: 'Czego kontrakt nie niesie', tresc: braki },
    ],
  });

  zapisz.addEventListener('click', () => {
    void zapiszRegule();
  });

  /** Zapis reguły na wskazanym układzie; stan bierze się z odpowiedzi. */
  async function zapiszRegule(): Promise<void> {
    const uklad = wiez.uklad();
    if (uklad === '') {
      powierzchnia.meldunek(
        'Harmonogram należy do UKŁADU automatyki, a układ nie jest wskazany. Wybierz go w podsekcji powiązania — reguła bez adresu nie ma czego uruchamiać.',
        false,
      );
      return;
    }
    const notacja = cykl.kontrolka.value.trim();
    const nazwaStrefy = strefa.kontrolka.value.trim();
    const wynik = await nadzor.zapiszHarmonogram({
      workflowId: uklad,
      enabled: obowiazuje.kontrolka.checked,
      ...(notacja === '' ? {} : { cron: notacja }),
      ...(nazwaStrefy === '' ? {} : { timeZone: nazwaStrefy }),
    });
    if (!wynik.udany || wynik.wynik === undefined) {
      powierzchnia.meldunek(`Rdzeń odmówił zapisu harmonogramu. ${powod(wynik.blad?.message, wynik.blad?.code)}`, false);
      return;
    }
    // Meldunek składa się z odpowiedzi rdzenia, łącznie z najbliższym uruchomieniem reguły harmonogramu.
    const zapisany = wynik.wynik;
    powierzchnia.meldunek(
      `Rdzeń zapisał harmonogram ${zapisany.id}: cron „${zapisany.cron ?? 'nie podany'}", ${zapisany.enabled ? 'obowiązuje' : 'nie obowiązuje'}${
        zapisany.nextRunAt === undefined
          ? '. Rdzeń nie podał najbliższego uruchomienia.'
          : `, najbliższe uruchomienie ${new Date(zapisany.nextRunAt).toLocaleString('pl-PL')}.`
      }`,
      true,
    );
    odczytaj();
  }

  /** Odczyt automatyk pod więź i harmonogramów wskazanego układu. */
  function odczytaj(): void {
    powierzchnia.tresci.ladowanie('Odczyt harmonogramów…');
    void (async () => {
      const automatyki = await nadzor.automatyki({});
      wiez.ustawWykaz(automatyki.udany ? (automatyki.wynik ?? []) : []);

      const uklad = wiez.uklad();
      const wynik = await nadzor.harmonogramy(uklad === '' ? {} : { workflowId: uklad });
      if (!wynik.udany || wynik.wynik === undefined) {
        powierzchnia.tresci.blad('Odczyt harmonogramów odmówiony.', wynik.blad);
        return;
      }
      powierzchnia.tresci.pusto('');
      if (wynik.wynik.length === 0) {
        lista.replaceChildren(
          zdaniePuste(
            uklad === ''
              ? 'Rdzeń nie ma zapisanego ani jednego harmonogramu. Praca ciągła nie jest stanem domyślnym — regułę zakłada Operator.'
              : 'Ten układ nie ma jeszcze reguły czasowej. Ustaw ją wyżej, a rdzeń zacznie uruchamiać układ sam.',
          ),
        );
        return;
      }
      lista.replaceChildren(...wynik.wynik.map(pozycjaHarmonogramu));
    })();
  }

  /** Jeden harmonogram jako pozycja wykazu wraz z jego czasami. */
  function pozycjaHarmonogramu(harmonogram: AutomationSchedule): HTMLElement {
    const { element } = pozycjaWykazu(
      harmonogram.cron ?? 'reguła bez notacji cron',
      `układ ${harmonogram.workflowId} · ${harmonogram.enabled ? 'obowiązuje' : 'nie obowiązuje'} · strefa ${harmonogram.timeZone ?? 'rdzenia'}`,
      PRZEDROSTEK,
    );
    element.append(
      wierszOpisu(
        'Najbliższe uruchomienie',
        harmonogram.nextRunAt === undefined
          ? 'rdzeń nie podał'
          : new Date(harmonogram.nextRunAt).toLocaleString('pl-PL'),
      ),
      wierszOpisu('Wyzwalacze poza cyklem', String(harmonogram.triggers?.length ?? 0)),
    );
    return element;
  }

  wiez.naZmiane(() => odczytaj());

  return {
    element: powierzchnia.element,
    odswiez: odczytaj,
    ustawOkno: powierzchnia.ustawOkno,
  };
}

/** Treść odmowy wraz z kodem kontraktu, złożona w jedno pełne zdanie gotowe do pokazania w meldunku sekcji. */
function powod(zdanie: string | undefined, kod: string | undefined): string {
  return `Powód: ${zdanie ?? 'rdzeń nie podał przyczyny'} (kod ${kod ?? 'brak'}).`;
}

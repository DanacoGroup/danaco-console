/**
 * Moduł MultitaskingAI składa scenę czterech okien ról pętli koordynator–wykonawca:
 * koordynator planuje bieg, dwa okna wykonawców pracują równolegle, analityk zestawia
 * wyniki. Wszystkie stoją na jednym stanie i jednym źródle biegu.
 */
import type { Kanal } from '../../protokol/kanal';
import { utworzZrodloSekcjiPaneli } from '../../powloka/zrodlo-sekcji-paneli';
import { widokZSesji, type OpisModulu } from '../rejestracja';
import { utworzWykazKomendRdzenia } from './braki-kontraktu';
import { utworzOknoAnalityka } from './okno-results-analyzer';
import { utworzPanelObsady } from './panel-obsady';
import { utworzPanelZadanWTle } from './panel-zadan-w-tle';
import { utworzZrodloPodagentow } from './zrodlo-podagentow';
import { utworzOknoKoordynatora } from './okno-coordinator-chat';
import { utworzOknoWykonawcy } from './okno-executor-chat';
import { utworzStanMultitaskingu } from './stan-multitaskingu';
import { utworzZrodloBiegu } from './zrodlo-biegu';
import { utworzZrodloOkien } from './zrodlo-okien';

/**
 * Liczba okien wykonawców obsadzanych na scenie. Wykaz środowiska w katalogu rdzenia
 * wymienia dwa okna Executor Chat, więc scena tworzy dokładnie tyle wystąpień okna
 * wykonawcy, a każde z nich otrzymuje własny numer porządkowy.
 */
const LICZBA_WYKONAWCOW = 2;

/**
 * Składa scenę czterech okien ról dla wskazanej sesji.
 *
 * Sesja jest tu wymagana, nie opcjonalna: role przypisuje się do OKIEN sesji,
 * więc bez niej nie ma czego obsadzić. Powłoka podaje ją przy wczytaniu.
 */
export function zamontujMultitasking(gospodarz: HTMLElement, kanal: Kanal, sesja: string): {
  zamknij(): void;
} {
  const bieg = utworzZrodloBiegu(kanal);
  const okna = utworzZrodloOkien(kanal);
  const podagenci = utworzZrodloPodagentow(kanal);
  const stan = utworzStanMultitaskingu(bieg, { sesja });
  // Wykaz komend rdzenia powstaje raz na scenę i jest wspólny czterem oknom.
  const komendy = utworzWykazKomendRdzenia(kanal);

  const scena = document.createElement('div');
  scena.className = 'dm-modul';
  scena.dataset['modul'] = 'multitasking';

  // Panel obsady stoi przed oknami ról i jest jedynym miejscem nadania tych ról.
  const sekcje = utworzZrodloSekcjiPaneli(kanal);
  const obsada = utworzPanelObsady({ zrodlo: okna, stan, sekcje });

  const koordynator = utworzOknoKoordynatora({ okna, bieg, stan, komendy });
  // Wykonawcy dostają to samo źródło podagentów co panel zadań w tle.
  const wykonawcy = Array.from({ length: LICZBA_WYKONAWCOW }, (_, i) =>
    utworzOknoWykonawcy({ okna, bieg, stan, podagenci, komendy, numer: i + 1 }),
  );
  const analityk = utworzOknoAnalityka({ okna, bieg, stan, komendy });

  // Panel zadań w tle jest jeden na kartę sesji i stoi pod oknami zlecającymi pracę.
  const zadaniaWTle = utworzPanelZadanWTle({ podagenci, bieg, stan });

  scena.append(
    obsada.element,
    koordynator.element,
    ...wykonawcy.map((w) => w.element),
    analityk.element,
    zadaniaWTle.element,
  );
  gospodarz.append(scena);

  // Pytanie o wykaz komend wychodzi przed pierwszym odświeżeniem okien.
  void komendy.odczytaj();
  obsada.odswiez();
  koordynator.odswiez();
  for (const wykonawca of wykonawcy) wykonawca.odswiez();
  analityk.odswiez();
  zadaniaWTle.odswiez();

  return {
    zamknij() {
      // Nasłuchy odpina się przed zdjęciem sceny, aby nie odświeżały zdjętego widoku.
      for (const wykonawca of wykonawcy) wykonawca.rozlacz();
      scena.remove();
    },
  };
}

/**
 * Samoopisujący się moduł dla rejestru powłoki.
 *
 * Kod `multitaskingai` odpowiada wierszowi środowiska w katalogu rdzenia —
 * nie literałowi wymyślonemu w kliencie.
 */
export const MODUL: OpisModulu = {
  kod: 'multitaskingai',
  utworzWidok: (kanal) =>
    widokZSesji(kanal, (gospodarz: HTMLElement, k: Kanal, sesja: string) =>
      zamontujMultitasking(gospodarz, k, sesja)),
};

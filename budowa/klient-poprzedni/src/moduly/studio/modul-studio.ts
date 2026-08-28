import './studio.css';

import type { Kanal } from '../../protokol/kanal';
import { utworzAparatZrodlo } from './aparat-zrodlo';
import { utworzObiektZrodlo } from './obiekt-zrodlo';
import { utworzOknoIngestOcrPanel } from './okno-ingest-ocr-panel';
import type { ZrodlaPostaciStudia } from './okno-pracy-z-dokumentem';
import { utworzSzablonZrodlo } from './szablon-zrodlo';
import { utworzTabeleZrodlo } from './tabela-zrodlo';
import { utworzZrodloKontroliStudio } from './zrodlo-kontroli-studio';
import { utworzZrodloPostaciStudio } from './zrodlo-postaci-studio';
import { utworzZrodloWstawienStudio } from './zrodlo-wstawien-studio';
import { utworzOknoPetliWykonawczej } from './okno-petli-wykonawczej';
import { utworzOknoPracyZDokumentem } from './okno-pracy-z-dokumentem';
import { utworzOknoRedakcjiDokumentu } from './okno-redakcji-dokumentu';
import { utworzOknoSessionRepository } from './okno-session-repository';
import { utworzOknoToolsPanel } from './okno-tools-panel';
import { utworzOknoWarsztatuDokumentu } from './okno-warsztatu-dokumentu';
import { utworzOsadzenieModulu } from './osadzenie-modulu';
import { utworzZrodloPetli } from './petla-zrodlo';
import { utworzZapleczePrzybornika } from './przybornik-zaplecze';
import { utworzStanStudio } from './stan-studio';
import { utworzZrodloAkcjiStudio } from './zrodlo-akcji-studio';
import { utworzZrodloDokumentuStudio } from './zrodlo-dokumentu-studio';
import { utworzZrodloMaterialuStudio } from './zrodlo-materialu-studio';
import { utworzZrodloOsadzenia } from './zrodlo-osadzenia';
import { utworzZrodloPracyStudio } from './zrodlo-pracy-studio';
import { utworzZrodloPrzekazania } from './zrodlo-przekazania';
import { utworzZrodloWarsztatuDokumentu } from './zrodlo-warsztatu-dokumentu';

/** Moduł Studio łączy jedno okno pracy z dokumentem i pięć okien bez powierzchni tekstowej, dawniej rozdzielonych na sześć okien. */
export interface ModulStudio {
  /** Element osadzany w obszarze roboczym powłoki. */
  element: HTMLElement;
  /** Zleca odczyt osadzenia, rejestru akcji i historii wersji. */
  wczytaj(idSesji: string): Promise<void>;
  /** Odłącza subskrypcję zdarzeń rdzenia. */
  rozlacz(): void;
}

export function utworzModulStudio(kanal: Kanal): ModulStudio {
  const stan = utworzStanStudio(kanal);
  const akcje = utworzZrodloAkcjiStudio(kanal);
  const przekazanie = utworzZrodloPrzekazania(kanal);
  const dokumenty = utworzZrodloDokumentuStudio(kanal);
  const material = utworzZrodloMaterialuStudio(kanal);
  const osadzenie = utworzOsadzenieModulu(stan, utworzZrodloOsadzenia(kanal));

  // Zaplecze przybornika stoi na poziomie modułu, bo jego nastawa dotyczy dwóch okien naraz.
  const przybornikZaplecze = utworzZapleczePrzybornika(kanal);

  // Siedem źródeł postaci stoi na poziomie modułu, bo cyfryzacja i okno pracy wołają to samo źródło.
  const zrodlaPostaci: ZrodlaPostaciStudia = {
    postaci: utworzZrodloPostaciStudio(kanal),
    kontrola: utworzZrodloKontroliStudio(kanal),
    wstawienia: utworzZrodloWstawienStudio(kanal),
    tabele: utworzTabeleZrodlo(kanal),
    obiekty: utworzObiektZrodlo(kanal),
    aparat: utworzAparatZrodlo(kanal),
    szablony: utworzSzablonZrodlo(kanal),
  };

  const cyfryzacja = utworzOknoIngestOcrPanel(
    stan,
    dokumenty,
    material,
    zrodlaPostaci.wstawienia,
  );
  const narzedzia = utworzOknoToolsPanel(stan, akcje, {
    naTryb: (tryb) => praca.ustawTrybOperacji(tryb),
  });
  const repozytorium = utworzOknoSessionRepository(stan);
  const zrodloWarsztatu = utworzZrodloWarsztatuDokumentu(kanal);
  const warsztat = utworzOknoWarsztatuDokumentu(stan, zrodloWarsztatu);
  const praca = utworzOknoPracyZDokumentem(
    stan,
    utworzZrodloPracyStudio(kanal),
    akcje,
    dokumenty,
    przekazanie,
    przybornikZaplecze,
    zrodlaPostaci,
    {
      naWiecejOperacji: () => narzedzia.przenieOgnisko(),
      naWarsztat: () => warsztat.przenieOgnisko(),
      naTrybOperacji: (tryb) => narzedzia.ustawTryb(tryb),
    },
  );
  // Redakcja bierze to samo źródło co warsztat: obie drogi wołają rdzeń tak samo.
  const redakcja = utworzOknoRedakcjiDokumentu(stan, zrodloWarsztatu);

  // Okno pętli wykonawczej wchodzi nakładką na żądanie i schodzi po zwinięciu, bez zabierania miejsca.
  const petla = utworzOknoPetliWykonawczej(stan, utworzZrodloPetli(kanal), {
    // Kliknięcie zadania prowadzi do miejsca, którego dotyczy: zakres staje się zaznaczeniem dokumentu.
    pokazWynikZadania: (zadanie) => {
      if (zadanie.rangeStart !== undefined && zadanie.rangeEnd !== undefined) {
        stan.ustawZaznaczenie({ poczatek: zadanie.rangeStart, koniec: zadanie.rangeEnd });
      }
      praca.element.scrollIntoView({ block: 'nearest' });
    },
    pokazFragment: (od, do_) => {
      stan.ustawZaznaczenie({ poczatek: od, koniec: do_ });
      praca.element.scrollIntoView({ block: 'nearest' });
    },
  });

  const pasWiodacy = document.createElement('div');
  pasWiodacy.className = 'ms-modul__pas ms-modul__pas--wiodacy';
  pasWiodacy.append(praca.element, narzedzia.element);

  // Pas pętli stoi nad pasem wiodącym, bo niesie jedyny znacznik prowadzący do rozwinięcia okna.
  const pasPetli = document.createElement('div');
  pasPetli.className = 'ms-modul__pas ms-modul__pas--petla';
  pasPetli.append(petla.element);

  // Warsztat dokumentu pracuje na materiale wniesionym do okna, nie na treści edytora.
  const pasWarsztatu = document.createElement('div');
  pasWarsztatu.className = 'ms-modul__pas ms-modul__pas--warsztat';
  pasWarsztatu.append(warsztat.element, redakcja.element);

  const element = document.createElement('div');
  element.className = 'ms-modul';
  element.dataset['modul'] = 'studio';
  element.setAttribute('aria-label', 'Moduł Studio — okna operacyjne');
  element.append(
    osadzenie.element,
    cyfryzacja.element,
    pasPetli,
    pasWiodacy,
    pasWarsztatu,
    repozytorium.element,
  );

  function odswiezWszystkie(): void {
    cyfryzacja.odswiez();
    praca.odswiez();
    narzedzia.odswiez();
    repozytorium.odswiez();
    warsztat.odswiez();
    redakcja.odswiez();
    petla.odswiez();
  }

  const odsubskrybuj = stan.obserwuj(odswiezWszystkie);
  odswiezWszystkie();

  return {
    element,

    async wczytaj(idSesji) {
      // Osadzenie idzie pierwsze, bo dopiero ono ustala identyfikator okna wymagany przez komendy modułu.
      await osadzenie.wczytaj(idSesji);
      await Promise.all([
        narzedzia.wczytaj(),
        repozytorium.wczytaj(),
        warsztat.wczytaj(),
        redakcja.wczytaj(),
        praca.wczytaj(),
        petla.wczytaj(),
      ]);
    },

    rozlacz() {
      odsubskrybuj();
      // Tools Panel zakłada nasłuch dokumentu przy menu z biblioteki, więc ma tu swoje osobne zamknięcie.
      narzedzia.zamknij();
      // Okno pętli trzyma osobną subskrypcję postępu, więc też wymaga zamknięcia przy rozbiórce.
      petla.zamknij();
      stan.rozlacz();
    },
  };
}

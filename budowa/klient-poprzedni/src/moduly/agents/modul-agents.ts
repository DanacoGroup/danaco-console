import './agents.css';

import type { Kanal } from '../../protokol/kanal';
import { utworzWykazBrakow, type WykazBrakow } from './braki-modulu';
import { utworzOknoAgentBuilder, type OknoAgentBuilder } from './okno-agent-builder';
import {
  utworzOknoConnectorsManager,
  type OknoConnectorsManager,
} from './okno-connectors-manager';
import {
  utworzOknoKatalogRozszerzen,
  type OknoKatalogRozszerzen,
} from './okno-katalog-rozszerzen';
import { utworzOknoModelConfiguration, type OknoModelConfiguration } from './okno-model-configuration';
import { utworzOknoPermissionsCenter, type OknoPermissionsCenter } from './okno-permissions-center';
import { utworzOknoSkillsManager, type OknoSkillsManager } from './okno-skills-manager';
import { utworzStanAgentow, type StanAgentow } from './stan-agentow';
import { utworzZrodloRozszerzen, type ZrodloRozszerzen } from './zrodlo-rozszerzen';
import { utworzZrodloZaplecza, type ZrodloZaplecza } from './zrodlo-zaplecza';

/**
 * Moduł Agents — sześć okien operacyjnych osadzonych w jednym układzie.
 *
 * Układ wynika z roli okna. Agent Builder jest kreatorem i punktem wejścia
 * modułu, więc stoi w pasie pierwszym na całą szerokość. Model Configuration
 * jest oknem pomocniczym, a cztery pozostałe są zarządcami, więc stoją w pasie
 * drugim obok siebie. Trzy z nich odnoszą się do eksperta wybranego w kreatorze;
 * Katalog rozszerzeń stoi na końcu pasa, bo jako jedyny mówi o platformie,
 * nie o ekspercie — i dlatego nie gaśnie, gdy żaden ekspert nie jest wybrany.
 *
 * Jeden ekspert jest czynny na cały moduł: `stan-agentow` jest jeden, więc wybór
 * w bibliotece przestawia okna eksperta naraz. Dlatego zakładka paska edytora
 * nie otwiera drugiego formularza — przenosi ognisko do okna, które daną rzeczą
 * zarządza.
 */
export interface ModulAgents {
  /** Element osadzany w obszarze roboczym powłoki. */
  element: HTMLElement;
  /** Zleca odczyt biblioteki, rejestrów i katalogów rdzenia. */
  wczytaj(idSesji: string): Promise<void>;
  /** Odłącza subskrypcję zdarzeń rdzenia. */
  rozlacz(): void;
}

export function utworzModulAgents(kanal: Kanal): ModulAgents {
  const stan: StanAgentow = utworzStanAgentow(kanal);
  const zaplecze: ZrodloZaplecza = utworzZrodloZaplecza(kanal);

  const model: OknoModelConfiguration = utworzOknoModelConfiguration(stan, zaplecze, kanal);
  const umiejetnosci: OknoSkillsManager = utworzOknoSkillsManager(stan, kanal);
  const konektory: OknoConnectorsManager = utworzOknoConnectorsManager(stan, zaplecze, kanal, {
    naOkno: (kod) => przenieOgnisko(kod),
  });
  const uprawnienia: OknoPermissionsCenter = utworzOknoPermissionsCenter(stan, zaplecze, kanal);
  const rozszerzenia: ZrodloRozszerzen = utworzZrodloRozszerzen(kanal);
  const katalog: OknoKatalogRozszerzen = utworzOknoKatalogRozszerzen(rozszerzenia);

  const builder: OknoAgentBuilder = utworzOknoAgentBuilder(stan, zaplecze, kanal, {
    naZakladke: (kod) => przenieOgnisko(kod),
  });

  const okna: Record<string, HTMLElement> = {
    'agent-builder': builder.element,
    'model-configuration': model.element,
    'skills-manager': umiejetnosci.element,
    'connectors-manager': konektory.element,
    'permissions-center': uprawnienia.element,
    'katalog-rozszerzen': katalog.element,
  };

  /**
   * Przenosi ognisko do okna wskazanego kodem zakładki.
   *
   * Kod zakładki nie jest kodem okna: zakładka „Tożsamość” prowadzi do okna
   * `agent-builder`. Pętla czyszcząca porównuje się z kodem okna, nie z kodem
   * zakładki — inaczej skasowałaby znacznik postawiony dla okna docelowego.
   */
  function przenieOgnisko(kod: string): void {
    const kodOkna = kod === 'tozsamosc' ? 'agent-builder' : kod;
    const cel = okna[kodOkna];
    if (cel === undefined) return;
    for (const [inny, wezel] of Object.entries(okna)) {
      if (inny === kodOkna) wezel.dataset['ognisko'] = 'tak';
      else delete wezel.dataset['ognisko'];
    }
    cel.scrollIntoView({ block: 'nearest' });
  }

  const pasZarzadcow = document.createElement('div');
  pasZarzadcow.className = 'da-modul__pas da-modul__pas--zarzadcy';
  pasZarzadcow.append(
    model.element,
    umiejetnosci.element,
    konektory.element,
    uprawnienia.element,
    katalog.element,
  );

  // Wykaz braków stoi na końcu modułu, poza pasami okien: mówi o module jako
  // całości, a nie o żadnym pojedynczym oknie. Pozycje mają też własne
  // kontrolki tam, gdzie Operator ich szuka; tutaj stoi ich liczba.
  const braki: WykazBrakow = utworzWykazBrakow(kanal);

  const element = document.createElement('div');
  element.className = 'da-modul';
  element.dataset['modul'] = 'agents';
  element.setAttribute('aria-label', 'Moduł Agents — okna operacyjne');
  element.append(builder.element, pasZarzadcow, braki.element);

  const odsubskrybuj = stan.obserwuj(() => {
    builder.odswiez();
    model.odswiez();
    umiejetnosci.odswiez();
    konektory.odswiez();
    uprawnienia.odswiez();
  });

  builder.odswiez();
  model.odswiez();
  umiejetnosci.odswiez();
  konektory.odswiez();
  uprawnienia.odswiez();

  return {
    element,

    async wczytaj(idSesji) {
      // Odczyty idą równolegle: każdy dotyczy innego obszaru kontraktu, a żaden
      // nie warunkuje drugiego. Odmowa jednego zostaje w jego oknie i nie gasi
      // pozostałych.
      await Promise.all([
        stan.odswiez(),
        builder.wczytajKategorie(),
        model.wczytaj(),
        konektory.wczytaj(),
        uprawnienia.wczytajOkna(idSesji),
        katalog.wczytaj(),
        braki.wczytaj(),
      ]);
    },

    rozlacz() {
      odsubskrybuj();
      builder.rozlacz();
      konektory.zamknij();
      // Drzewo wyboru narzędzi zakłada nasłuchy na dokumencie (zamknięcie
      // kliknięciem obok, Escape), więc bez tego wiersza przeżyłoby własne okno
      // i reagowało na klawiaturę w module, którego już nie ma.
      umiejetnosci.zamknij();
      model.zamknij();
      uprawnienia.zamknij();
      braki.zamknij();
      stan.rozlacz();
    },
  };
}

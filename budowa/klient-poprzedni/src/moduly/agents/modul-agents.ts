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
 * Moduł Agents łączy sześć okien operacyjnych w jednym układzie: Agent Builder w pasie
 * pierwszym, cztery zarządcy i konfiguracja modelu w pasie drugim.
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

  /** Przenosi ognisko do okna wskazanego kodem zakładki; kod zakładki bywa inny niż kod okna. */
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

  // Wykaz braków stoi poza pasami okien — dotyczy modułu jako całości, nie pojedynczego okna.
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
      // Odczyty idą równolegle: dotyczą różnych obszarów kontraktu i nie warunkują się nawzajem.
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
      // Drzewo wyboru narzędzi zakłada nasłuchy na dokumencie, więc wymaga jawnego zamknięcia.
      umiejetnosci.zamknij();
      model.zamknij();
      uprawnienia.zamknij();
      braki.zamknij();
      stan.rozlacz();
    },
  };
}

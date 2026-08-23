// Źródła SVG — grupa: srodowiska.
// Nazwy i kolejność wprost z `ikony/manifest.json` (pozycje 79–82).

import srodowiskoTalkin from '../svg/srodowisko-talkin.svg?raw';
import srodowiskoWorkspace from '../svg/srodowisko-workspace.svg?raw';
import srodowiskoCodestudio from '../svg/srodowisko-codestudio.svg?raw';
import srodowiskoMultitaskingai from '../svg/srodowisko-multitaskingai.svg?raw';

/** Emblematy czterech środowisk — znaki własne, nie Lucide. */
export const SRODOWISKA = {
  'srodowisko-talkin': srodowiskoTalkin,
  'srodowisko-workspace': srodowiskoWorkspace,
  'srodowisko-codestudio': srodowiskoCodestudio,
  'srodowisko-multitaskingai': srodowiskoMultitaskingai,
} as const;

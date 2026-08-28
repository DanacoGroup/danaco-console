import '../motyw/motyw.css';
import '../komponenty/indeks.css';

import { ProgressStatus } from '../../../shared/contract';
import { uruchomMotyw, zastosujMotyw } from '../motyw/motyw';
import { utworzPowloke } from './powloka';

// Podgląd powłoki środowiska poza aplikacją — strona sprawdzająca cztery pasy przed ich osadzeniem.

uruchomMotyw();

// Oba motywy są równoprawne, więc podgląd musi umieć pokazać każdy z nich
// bez ruszania ustawień środowiska.
const zadany = new URLSearchParams(window.location.search).get('motyw');
if (zadany === 'jasny') zastosujMotyw('light');
if (zadany === 'ciemny') zastosujMotyw('dark');

const powloka = utworzPowloke({ operator: 'Dariusz Naharnowicz' });

// Trzy dodatkowe karty sesji pokazują wskaźnik pracy w tle w każdym z jego trzech możliwych stanów pracy.
powloka.karty.dodaj('Analiza rynku', ProgressStatus.Running);
powloka.karty.dodaj('Umowa ramowa', ProgressStatus.Paused);
powloka.karty.dodaj('Research OZE', ProgressStatus.Done);
powloka.pasek.ustawPowiadomienia(3);

const pierwsza = powloka.karty.wykaz()[0];
if (pierwsza !== undefined) powloka.karty.wybierz(pierwsza.id);

document.body.append(powloka.element);

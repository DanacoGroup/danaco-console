// Plik nie jest osiągalny z `main.ts` — jedyna droga do niego prowadzi przez
// własny punkt wejścia Vite `podglad-ukladu.html`.
import '../aplikacja/arkusze-stylow';

import { uruchomMotyw } from '../motyw/motyw';
import {
  odczytajParametry,
  zastosujMotywPodgladu,
  zastosujParametry,
} from './parametry-podgladu';
import { utworzPulpitProby } from './pulpit-proby';
import { zasiejTrescProbna } from './tresc-probna';
import { utworzUkladOkien } from './uklad-okien';

// Punkt wejścia strony podglądu układu, łączący motyw, treść przykładową i parametry adresu.

const gospodarz = document.querySelector<HTMLElement>('#danaco-podglad') ?? document.body;
const parametry = odczytajParametry(window.location.href);

uruchomMotyw();
zastosujMotywPodgladu(parametry);

// Stanowisko podglądu nie ma rdzenia i nie osadza rozmowy prowadzonej w rzeczywistej aplikacji, dlatego układ okien otrzymuje puste gniazda w każdym uruchomieniu podglądu.
const uklad = utworzUkladOkien({ wbudowanaRozmowa: true });

gospodarz.append(uklad.element, utworzPulpitProby(uklad));
zasiejTrescProbna(uklad);
zastosujParametry(uklad, parametry);

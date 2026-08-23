; Odpowiedzialność pliku: skróty, których szablon NSIS Tauriego sam nie zakłada.
;
; Szablon zakłada dwa skróty i oba wystarczają same z siebie: wpis w menu Start
; (`CreateOrUpdateStartMenuShortcut`, prosto w `$SMPROGRAMS`, bo `STARTMENUFOLDER`
; jest pusty) oraz skrót na pulpicie (`CreateOrUpdateDesktopShortcut`, spod
; zaznaczenia na stronie końcowej, a przy instalacji cichej i przelotnej —
; bezwarunkowo). Tego pliku nie ma po co dublować.
;
; Brakuje trzeciego miejsca, o które prosi Właściciel: szybkiego wybierania.
; Zakłada je ten plik, bo Tauri nie ma na nie ani nastawy, ani makra.
;
; PRZYPIĘCIA DO PASKA ZADAŃ NIE MA I NIE BĘDZIE — nie z braku starania.
; Windows od wersji 10 nie wystawia czasownika `taskbarpin` do wywołania
; z programu: powłoka odmawia go każdemu, kto nie jest użytkownikiem klikającym
; w menu. Nie da się tego obejść ani z NSIS, ani z powłoki natywnej; jedyne
; wsparte drogi to zasady grupy (`LayoutModification.xml`, wdrożenie firmowe)
; albo ręka użytkownika. Sam Tauri mówi to samo swoim zestawem makr: `utils.nsh`
; ma `UnpinShortcut` — odpiąć przy odinstalowaniu potrafi — i nie ma ani jednego
; makra przypinającego.
;
; Co robimy zamiast obietnicy: pilnujemy, żeby przypięcie ręczne DZIAŁAŁO
; POPRAWNIE. Każdy zakładany tu skrót dostaje `AppUserModelId` tym samym makrem,
; którego szablon używa dla menu Start i pulpitu. Bez tego znacznika Windows
; uznaje okno powłoki i przypięty skrót za dwa różne programy — użytkownik
; przypina ikonę, a po uruchomieniu dostaje na pasku drugą, obok pierwszej.
;
; Wpięcie: `bundle.windows.nsis.installerHooks` w profilu Tauriego.

!macro NSIS_HOOK_POSTINSTALL
  ; Aktualizacja i `/NS` nie zakładają skrótów — tak samo, jak w szablonie.
  ; Bez tego warunku aktualizacja wskrzeszałaby skrót, który użytkownik usunął.
  ${If} $UpdateMode = 0
  ${AndIf} $NoShortcutMode = 0
    CreateShortcut "$QUICKLAUNCH\${PRODUCTNAME}.lnk" "$INSTDIR\${MAINBINARYNAME}.exe"
    !insertmacro SetLnkAppUserModelId "$QUICKLAUNCH\${PRODUCTNAME}.lnk"
    DetailPrint "Skrót w szybkim wybieraniu: $QUICKLAUNCH\${PRODUCTNAME}.lnk"
  ${EndIf}
!macroend

!macro NSIS_HOOK_POSTUNINSTALL
  ; Odpięcie przed usunięciem: pasek zadań trzyma własną kopię skrótu i sam plik
  ; usunięty spod niej zostawia ikonę-widmo, która po kliknięciu nie robi nic.
  !insertmacro UnpinShortcut "$QUICKLAUNCH\${PRODUCTNAME}.lnk"
  Delete "$QUICKLAUNCH\${PRODUCTNAME}.lnk"
!macroend

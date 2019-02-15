set nocompatible
set number
set guioptions-=r
set guioptions-=L
set guioptions-=b

noremap <leader>y "+y
nmap <leader>p "+p

set showtabline=1
set guifont=Fira_Code_Retina:h12,DejaVu_Sans_Mono:h12
set guifont=Fira\ Code\ Retina\ 12,DejaVu\ Sans\ Mono\ 12

if has('gui_running')
    set guioptions-=T  " no toolbar
    if has('gui_win32')
      set guifont=Fira_Code_Retina:h12,DejaVu_Sans_Mono:h12:cANSI
    else
      set guifont=Fira\ Code\ Retina\ 12,DejaVu\ Sans\ Mono\ 12
    endif
  endif

syntax on
syntax enable
set background=dark
set nowrap
set fileformat=unix
set fileformats=unix,dos,mac
set cindent
set autoindent
filetype indent on
filetype plugin on
set expandtab
set tabstop=4
set shiftwidth=4
set softtabstop=4
set history=200

set showmatch
set scrolloff=5
set laststatus=2
set fileencodings=utf-8,gb2312,gb18030,gbk,ucs-bom,cp936,latin1
set enc=utf8
set fencs=utf8,gbk,gb2312,gb18030
set encoding=utf-8
set backspace=2
set linespace=6
set mouse=a
set selection=exclusive
set selectmode=mouse,key
set clipboard=unnamed
set matchtime=5
set ignorecase
set incsearch
set hlsearch
set noexpandtab
set whichwrap+=<,>,h,l
set autoread
set cursorline
set cursorcolumn
set nu
set cmdheight=5

nnoremap <silent> [b :bprevious<CR>
nnoremap <silent> ]b :bnext<CR>
nnoremap <silent> [B :bfirst<CR>
nnoremap <silent> ]B :blast<CR>
" remap U to <C-r>"
nnoremap U <C-r>
"clear highlights on hitting esc twice."
nnoremap <esc><esc> :noh<return>

call plug#begin('~/.vim/plugged')

Plug 'scrooloose/nerdtree'
Plug 'nathanaelkane/vim-indent-guides'
Plug 'scrooloose/nerdcommenter'
Plug 'wakatime/vim-wakatime'
Plug 'octol/vim-cpp-enhanced-highlight'
Plug 'jiangmiao/auto-pairs'
Plug 'Xuyuanp/nerdtree-git-plugin'
Plug 'ludovicchabant/vim-gutentags'
Plug 'Yggdroot/LeaderF'
Plug 'skywind3000/asyncrun.vim' , { 'as': 'asyncrun' }
Plug 'Vimjas/vim-python-pep8-indent'
Plug 'tweekmonster/impsort.vim'
Plug 'kien/ctrlp.vim'
Plug 'ntpeters/vim-better-whitespace'
Plug 'dracula/vim', { 'as': 'dracula' }
Plug 'Art2Cat/vim-autoformat'
Plug 'fatih/vim-go'
Plug 'majutsushi/tagbar'
Plug 'vim-airline/vim-airline'
Plug 'airblade/vim-gitgutter'
Plug 'vim-syntastic/syntastic'
Plug 'tpope/vim-commetary'
Plug 'w0rp/ale'
Plug 'davidhalter/jedi-vim'
call plug#end()

filetype on

colorscheme dracula

"NERDTree configuration""
nnoremap <F2> :NERDTreeToggle<CR>

let NERDTreeChDirMode=1
let NERDTreeShowBookmarks=1
let NERDTreeIgnore=['\~$', '\.pyc$', '\.swp$']
let NERDTreeWinSize=25
let NERDTreeShowFiles=1
let NERDTreeShowLineNumbers=1
let NERDTreeShowHidden=0
let NERDTreeHightCursorline=1
let NERDTreeAutoCenter=1
let NERDTreeChristmasTree=1
let NERDTreeMapOpenInTab='<ENTER>'

let g:NERDTreeIndicatorMapCustom={
			\ "Modified"  : "✹",
			\ "Staged"    : "✚",
			\ "Untracked" : "✭",
			\ "Renamed"   : "➜",
			\ "Unmerged"  : "═",
			\ "Deleted"   : "✖",
			\ "Dirty"     : "✗",
			\ "Clean"     : "✔︎",
			\ "Unknown"   : "?"
			\ }


let mapleader=','

let g:gutentags_project_root=['.root', '.svn', '.git', '.hg', '.project']
let g:gutentags_ctags_tagfile='.tags'
let s:vim_tags=expand('~/.cache/tags')
let g:gutentags_cache_dir=s:vim_tags
let g:gutentags_ctags_extra_args=['--fields=+niazS', '--extra=+q']
let g:gutentags_ctags_extra_args += ['--c++-kinds=+px']
let g:gutentags_ctags_extra_args += ['--c-kinds=+px']
let g:gutentags_ctags_extra_args += ['--rust-kinds=+px']


if !isdirectory(s:vim_tags)
	    silent! call mkdir(s:vim_tags, 'p')
endif

" asyncrun configuration

function! s:compile_and_run()
	exec 'w'
	exec 'vertical rightbelow copen 80'
	exec 'wincmd w'
	if &filetype ==# 'c'
		exec 'AsyncRun! gcc -Wall % -o %<; time ./%<'
	elseif &filetype ==# 'cpp'
		exec 'AsyncRun! g++ -Wall -std=c++11 % -o %<; time ./%<'
	elseif &filetype ==# 'rust'
		exec 'AsyncRun! rustc %; time ./%<'
	elseif &filetype ==# 'java'
		exec 'AsyncRun! javac %; time java %<; rm -f *.class'
	elseif &filetype ==# 'sh'
		exec 'AsyncRun! time bash %'
	elseif &filetype ==# 'python'
		exec 'AsyncRun! time python3 "%"'
	elseif &filetype ==# 'javascript'
		exec 'AsyncRun! time node %'
	elseif &filetype ==# 'go'
		exec 'AsyncRun! time go run %'
	endif
endfunction
nnoremap <F5> :call <SID>compile_and_run()<CR>
" let g:asyncrun_open=6
nnoremap <F10> :call asyncrun#quickfix_toggle(6)<cr>
nnoremap <slient> <F9> :AsyncRun g++ -Wall -02 "$(VIM_FILEPATH)" -o "$(VIM_FILEDIR)/$(VIM_FILENOEXT)" <cr>
nnoremap <slient> <F8> :AsyncRun -raw -cwd=$(VIM_FILEPATH) $(VIM_FILEDIR)/$(VIM_FILENOEXT)" <cr>

let g:better_whitespace_enable=1
let g:strip_whitespace_on_save=1

noremap <F3> :Autoformat<CR>

let g:indent_guides_enable_on_vim_startup=1

if has("gui_running")
	" GUI is running or is about to start.
	" Maximize gvim window (for an alternative on Windows, see simalt below).
	set lines=999 columns=999
else
	" This is console Vim.
	if exists("+lines")
		set lines=50
	endif
	if exists("+columns")
		set columns=100
	endif
endif

" Disable Arrow keys in Escape mode
map <up> <nop>
map <down> <nop>
map <left> <nop>
map <right> <nop>

" Disable Arrow keys in Insert mode
" imap <up> <nop>
" imap <down> <nop>
" imap <left> <nop>
" imap <right> <nop>

nmap <F7> :TagbarToggle<CR>

let g:airline#extensions#tabline#enabled = 1


let g:go_fmt_command = "goimports"
let g:go_addtags_transform = "camelcase"
let g:go_highlight_types = 1
let g:go_highlight_fields = 1
let g:go_highlight_functions = 1
let g:go_highlight_function_calls = 1
autocmd BufNewFile,BufRead *.go setlocal noexpandtab tabstop=4 shiftwidth=4

let g:python3_host_prog = '/usr/bin/python3'

" ale options
let g:ale_python_flake8_options = '--ignore=E129,E501,E302,E265,E241,E305,E402,W503'
let g:ale_python_pylint_options = '-j 0 --max-line-length=120'
let g:ale_list_window_size = 4
let g:ale_sign_column_always = 0
let g:ale_open_list = 1
let g:ale_keep_list_window_open = '1'

" Options are in .pylintrc!
highlight VertSplit ctermbg=253

let g:ale_sign_error = '‼'
let g:ale_sign_warning = '∙'
let g:ale_lint_on_text_changed = 'never'
let g:ale_lint_on_enter = '0'
let g:ale_lint_on_save = '1'
nmap <silent> <C-M> <Plug>(ale_previous_wrap)
nmap <silent> <C-m> <Plug>(ale_next_wrap)

" mapping to make movements operate on 1 screen line in wrap mode
function! ScreenMovement(movement)
   if &wrap
      return "g" . a:movement
   else
      return a:movement
   endif
endfunction

" highlight python and self function
autocmd BufEnter * syntax match Type /\v\.[a-zA-Z0-9_]+\ze(\[|\s|$|,|\]|\)|\.|:)/hs=s+1
autocmd BufEnter * syntax match pythonFunction /\v[[:alnum:]_]+\ze(\s?\()/
hi def link pythonFunction Function
autocmd BufEnter * syn match Self "\(\W\|^\)\@<=self\(\.\)\@="
highlight self ctermfg=239

" jedi options
let g:jedi#auto_initialization = 1
let g:jedi#completions_enabled = 0
let g:jedi#auto_vim_configuration = 0
let g:jedi#smart_auto_mappings = 0
let g:jedi#popup_on_dot = 0
let g:jedi#completions_command = ""
let g:jedi#show_call_signatures = "1"
let g:jedi#show_call_signatures_delay = 0
let g:jedi#use_tabs_not_buffers = 0
let g:jedi#show_call_signatures_modes = 'i'  " ni = also in normal mode
let g:jedi#enable_speed_debugging=0

" Impsort option
hi pythonImportedObject ctermfg=127
hi pythonImportedFuncDef ctermfg=127
hi pythonImportedClassDef ctermfg=127

" vimgutter options
let g:gitgutter_override_sign_column_highlight = 0
let g:gitgutter_map_keys = 0
 
" ctrl p options
let g:ctrlp_custom_ignore = '\v\.(npy|jpg|pyc|so|dll)$'
let g:ctrlp_user_command = ['.git', 'cd %s && git ls-files -co --exclude-standard']
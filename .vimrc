"no compatible with vi"
set nocompatible
"show line"
set number
" hide scroll"
set guioptions-=r
set guioptions-=L
set guioptions-=b

vnoremap <Leader>y "+y
nmap <Leader>p "+p

"hide top tag"
set showtabline=1
"set font"
set guifont=Menlo:h14,DejaVu_Sans_Mono:h14
set guifont=Menlo\ 14,DejaVu\ Sans\ Mono\ 14
syntax on   " enable syntax
syntax enable
set background=dark   "set background"
" set termguicolors
colorscheme dracula
color dracula
set nowrap  "set nowrap line"
set fileformat=unix "fileformat unix 'lf'"
" set fileformats=unix,dos,mac
set cindent     "indenct like C"
set autoindent
filetype indent on
filetype plugin on
set expandtab
set tabstop=2   "set tab size 4"
set shiftwidth=2
set softtabstop=2
set history=200

set showmatch   "show match parentheses"
set scrolloff=5     "5 rows from the top and bottom"
set laststatus=2    "Command line 2 lines"
set fileencodings=utf-8,gb2312,gb18030,gbk,ucs-bom,cp936,latin1
set enc=utf8
set fencs=utf8,gbk,gb2312,gb18030      "file encoding"
set encoding=utf-8
set backspace=2
set linespace=6
set mouse=a     "enable mouse"
set selection=exclusive
set selectmode=mouse,key
set clipboard=unnamed
set matchtime=5
set ignorecase      "ignore upper or lower case"
set incsearch
set hlsearch        "search high light"
set noexpandtab     "not allow expand table"
set whichwrap+=<,>,h,l
set autoread
set cursorline      " highlight current line"
set cursorcolumn        "highlight current column"
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

Plug 'Valloric/YouCompleteMe'
"Plug 'maralla/completor.vim' , { 'as': 'completor' }
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
Plug 'ntpeters/vim-better-whitespace'
Plug 'dracula/vim', { 'as': 'dracula' }
Plug 'Chiel92/vim-autoformat'
Plug 'fatih/vim-go'
Plug 'rust-lang/rust.vim', { 'as': 'rustvim' }
Plug 'majutsushi/tagbar'
Plug 'vim-airline/vim-airline'
Plug 'airblade/vim-gitgutter'
call plug#end()

filetype on

autocmd FileType python nnoremap <LocalLeader>i :!isort %<CR><CR>


"YouCompleteMe configuration"
let g:ycm_global_ycm_extra_conf='~/.vim/plugged/YouCompleteMe/third_party/ycmd/.ycm_extra_conf.py'
let g:ycm_rust_src_path = '$(rustc --print sysroot)/lib/rustlib/src/rust/src'
let g:ycm_confirm_extra_conf=1
set completeopt=longest,menu
let g:ycm_python_binary_path='python'
let g:ycm_seed_identifiers_with_syntax=1
let g:ycm_complete_in_comments=1
let g:ycm_collect_identifiers_from_tags_files=1
let g:ycm_collect_identifiers_from_comments_and_strings=0
let g:ycm_min_num_of_chars_for_completion=2
let g:ycm_autoclose_preview_window_after_completion=1
let g:ycm_cache_omnifunc=0
let g:ycm_complete_in_strings=1
let g:ycm_key_invoke_completion='<C-Space>'
let g:ycm_register_as_syntastic_checker=1
let g:show_diagnostics_ui=1
let g:ycm_enable_diagnostic_signs=1
let g:ycm_enable_diagnostic_highlighting=1
let g:ycm_always_populate_location_list=1
let g:ycm_semantic_triggers={
			\ 'c,cpp,python,java,go,erlang,perl,rust': ['re!\w{2}'],
			\ 'cs,lua,javascript': ['re!\w{2}'],
			\ }
nnoremap <leader>jd :YcmCompleter GoToDefinitionElseDeclaration<CR>

"NERDTree configuration""
nnoremap <F2> :NERDTreeToggle<CR>
"map <F2> :NERDTreeMirror<CR>":

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

" gutentags 搜索工程目录的标志，碰到这些文件/目录名就停止向上一级目录递归
let g:gutentags_project_root=['.root', '.svn', '.git', '.hg', '.project']

" 所生成的数据文件的名称
let g:gutentags_ctags_tagfile='.tags'

" 将自动生成的 tags 文件全部放入 ~/.cache/tags 目录中，避免污染工程目录
let s:vim_tags=expand('~/.cache/tags')
let g:gutentags_cache_dir=s:vim_tags

" 配置 ctags 的参数
let g:gutentags_ctags_extra_args=['--fields=+niazS', '--extra=+q']
let g:gutentags_ctags_extra_args += ['--c++-kinds=+px']
let g:gutentags_ctags_extra_args += ['--c-kinds=+px']
let g:gutentags_ctags_extra_args += ['--rust-kinds=+px']


" 检测 ~/.cache/tags 不存在就新建
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

" let g:completor_clang_binary='/usr/bin/clang'
"map <F3> <Plug>CompletorCppJumpToPlaceholder
"imap <F3> <Plug>CompletorCppJumpToPlaceholder

au BufWrite * :Autoformat

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

import { useEffect, useMemo, useState } from 'react'
import { useDispatch, useSelector } from 'react-redux'
import { useLocation, useNavigate } from 'react-router-dom'
import {
  AppBar,
  Box,
  CssBaseline,
  Drawer,
  IconButton,
  List,
  ListItem,
  ListItemButton,
  ListItemIcon,
  ListItemText,
  Toolbar,
  Typography,
  Avatar,
  Menu,
  MenuItem,
  Divider,
  Chip,
  Tooltip,
} from '@mui/material'
import {
  Menu as MenuIcon,
  Dashboard,
  TrendingUp,
  AccountBalance,
  Factory,
  People,
  Inventory,
  LocalShipping,
  Logout,
} from '@mui/icons-material'
import { logout } from '../../store/slices/authSlice'

const drawerWidth = 240

const menuItems = [
  { text: 'Executive BI', icon: <Dashboard />, path: '/bi', roles: ['executive', 'analyst'] },
  { text: 'Reports', icon: <Dashboard />, path: '/reports', roles: ['executive', 'analyst'] },
  { text: 'Integration Jobs', icon: <Dashboard />, path: '/integration', roles: ['executive'] },
  { text: 'Sales', icon: <TrendingUp />, path: '/sales', roles: ['executive', 'sales_manager', 'analyst'] },
  { text: 'Finance', icon: <AccountBalance />, path: '/finance', roles: ['executive', 'finance_manager', 'analyst'] },
  { text: 'Production', icon: <Factory />, path: '/production', roles: ['executive', 'production_manager', 'analyst'] },
  { text: 'Inventory', icon: <Inventory />, path: '/inventory', roles: ['executive', 'inventory_manager', 'analyst', 'supply_chain_manager'] },
  { text: 'HR', icon: <People />, path: '/hr', roles: ['executive', 'hr_manager', 'analyst'] },
  { text: 'Supply Chain', icon: <LocalShipping />, path: '/supplychain', roles: ['executive', 'supply_chain_manager', 'analyst'] },
]

function Layout({ children }) {
  const [mobileOpen, setMobileOpen] = useState(false)
  const [anchorEl, setAnchorEl] = useState(null)
  const dispatch = useDispatch()
  const navigate = useNavigate()
  const location = useLocation()
  const { user, roles, defaultCompanyId, companyIds } = useSelector((state) => state.auth)

  const lowerRoles = useMemo(() => roles?.map((role) => role.toLowerCase()) || [], [roles])

  const handleDrawerToggle = () => {
    setMobileOpen((prev) => !prev)
  }

  const handleMenuOpen = (event) => {
    setAnchorEl(event.currentTarget)
  }

  const handleMenuClose = () => {
    setAnchorEl(null)
  }

  const handleLogout = () => {
    dispatch(logout())
    navigate('/login')
  }

  const filteredMenuItems = useMemo(
    () =>
      menuItems.filter((item) =>
        item.roles.some((role) => lowerRoles.includes(role))
      ),
    [lowerRoles]
  )

  useEffect(() => {
    if (mobileOpen) {
      setMobileOpen(false)
    }
  }, [location.pathname])

  const drawer = (
    <Box>
      <Toolbar>
        <Box>
          <Typography variant="h6" noWrap component="div">
            iDash
          </Typography>
          {defaultCompanyId && (
            <Typography variant="caption" color="text.secondary">
              Company #{defaultCompanyId}
            </Typography>
          )}
        </Box>
      </Toolbar>
      <Divider />
      <List>
        {filteredMenuItems.map((item) => {
          const isActive = location.pathname.startsWith(item.path)
          return (
            <ListItem key={item.text} disablePadding>
              <Tooltip title={item.text} placement="right" enterDelay={400} disableInteractive>
                <ListItemButton
                  onClick={() => navigate(item.path)}
                  selected={isActive}
                  aria-current={isActive ? 'page' : undefined}
                  sx={{
                    borderLeft: isActive ? (theme) => `4px solid ${theme.palette.primary.main}` : '4px solid transparent',
                  }}
                >
                  <ListItemIcon sx={{ minWidth: 36 }}>{item.icon}</ListItemIcon>
                  <ListItemText primary={item.text} />
                </ListItemButton>
              </Tooltip>
            </ListItem>
          )
        })}
      </List>
      {companyIds?.length > 1 && (
        <Box sx={{ px: 2, py: 1 }}>
          <Typography variant="caption" color="text.secondary">
            Linked Companies
          </Typography>
          <Box display="flex" flexWrap="wrap" gap={0.5} mt={1}>
            {companyIds.map((id) => (
              <Chip key={id} label={`#${id}`} size="small" color={id === defaultCompanyId ? 'primary' : 'default'} />
            ))}
          </Box>
        </Box>
      )}
    </Box>
  )

  return (
    <Box sx={{ display: 'flex' }}>
      <CssBaseline />
      <AppBar
        position="fixed"
        sx={{
          bgcolor: 'background.paper',
          color: 'text.primary',
          boxShadow: (theme) => theme.shadows[1],
          width: { sm: `calc(100% - ${drawerWidth}px)` },
          ml: { sm: `${drawerWidth}px` },
        }}
      >
        <Toolbar>
          <IconButton
            color="inherit"
            aria-label="open navigation"
            edge="start"
            onClick={handleDrawerToggle}
            sx={{ mr: 2, display: { sm: 'none' } }}
          >
            <MenuIcon />
          </IconButton>
          <Typography variant="h6" noWrap component="div" sx={{ flexGrow: 1 }}>
            Intelligent Management Dashboard
          </Typography>
          <Box display="flex" alignItems="center" gap={1.5}>
            {lowerRoles.map((role) => (
              <Chip key={role} label={role.replace('_', ' ')} size="small" color="secondary" sx={{ textTransform: 'capitalize' }} />
            ))}
            <IconButton onClick={handleMenuOpen} color="inherit" aria-label="Account menu">
              <Avatar sx={{ width: 32, height: 32 }}>
                {user?.name?.charAt(0).toUpperCase()}
              </Avatar>
            </IconButton>
          </Box>
          <Menu
            anchorEl={anchorEl}
            open={Boolean(anchorEl)}
            onClose={handleMenuClose}
            anchorOrigin={{ vertical: 'bottom', horizontal: 'right' }}
            transformOrigin={{ vertical: 'top', horizontal: 'right' }}
          >
            <MenuItem onClick={handleLogout}>
              <ListItemIcon>
                <Logout fontSize="small" />
              </ListItemIcon>
              Logout
            </MenuItem>
          </Menu>
        </Toolbar>
      </AppBar>
      <Box
        component="nav"
        sx={{ width: { sm: drawerWidth }, flexShrink: { sm: 0 } }}
      >
        <Drawer
          variant="temporary"
          open={mobileOpen}
          onClose={handleDrawerToggle}
          ModalProps={{
            keepMounted: true,
          }}
          sx={{
            display: { xs: 'block', sm: 'none' },
            '& .MuiDrawer-paper': {
              boxSizing: 'border-box',
              width: drawerWidth,
            },
          }}
        >
          {drawer}
        </Drawer>
        <Drawer
          variant="permanent"
          sx={{
            display: { xs: 'none', sm: 'block' },
            '& .MuiDrawer-paper': {
              boxSizing: 'border-box',
              width: drawerWidth,
            },
          }}
          open
        >
          {drawer}
        </Drawer>
      </Box>
      <Box
        component="main"
        sx={{
          flexGrow: 1,
          minHeight: '100vh',
          bgcolor: 'background.default',
          p: { xs: 2, sm: 3 },
          width: { sm: `calc(100% - ${drawerWidth}px)` },
        }}
      >
        <Toolbar />
        {children}
      </Box>
    </Box>
  )
}

export default Layout

